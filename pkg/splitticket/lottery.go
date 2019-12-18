package splitticket

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/crypto/blake256"
	"github.com/decred/dcrd/dcrutil"
)

const (
	// SecretNbHashSize is the size of the hash of a secret number
	SecretNbHashSize = 32

	// SecretNbSize is the number of random bytes of a secret number
	SecretNbSize = 8

	// LotteryCommitmentHashSize is the size in bytes of the lottery commitment
	LotteryCommitmentHashSize = 32
)

// SelectContributionAmounts decides how to split the ticket priced at
// ticketPrice such that each ith participant contributes at most maxAmount[i],
// pays a participation fee partFee and everyone pays total poolFee.
//
// In order to call this function sum(maxAmounts) must be > (ticketPrice +
// sum(partFee)) and the maxAmounts **MUST** be ordered in ascending order of
// amount.
//
// The algorithm tries to split the purchasing amounts between the participants
// in such a way as to average the contribution % of the participants, while
// ensuring that the ticket can be bought.
//
// This function returns the corresponding commitment amounts for each
// participant (ie, what the sstxcommitment output will be for each
// participant) and the proportional pool fee amount.
func SelectContributionAmounts(maxAmounts []dcrutil.Amount, ticketPrice,
	partFee, poolFee dcrutil.Amount) ([]dcrutil.Amount, []dcrutil.Amount, error) {

	nparts := len(maxAmounts)
	totalAvailable := dcrutil.Amount(0)
	prev := dcrutil.Amount(0)
	remainingMax := make([]dcrutil.Amount, nparts)
	for i, a := range maxAmounts {
		if a < partFee {
			return nil, nil, fmt.Errorf("Invalid Argument: amount in index %d (%s)"+
				"less than minimum needed (%s)", i, a, partFee)
		}

		if a < 0 {
			return nil, nil, fmt.Errorf("Invalid Argument: amount is negative (%s)"+
				"in index %d", a, i)
		}

		if a < prev {
			return nil, nil, fmt.Errorf("Invalid Argument: maxAmounts must be in " +
				"ascending order")
		}
		totalAvailable += a
		prev = a

		// remainingMax considers the partFee paied by the participant
		remainingMax[i] = a - partFee
	}

	// notice that while the partFee is an "additional" amount needed by each
	// participant, the pool fee is simply a part of the ticket price placed
	// into the first sstx output, therefore it is already accounted as needed
	// in the ticketPrice amount
	neededAmount := ticketPrice + partFee*dcrutil.Amount(nparts)
	if totalAvailable < neededAmount {
		return nil, nil, fmt.Errorf("Invalid Argument: total available %s less "+
			"than needed %s", totalAvailable, neededAmount)
	}

	// totalLeft is the remaining part of the ticket that needs to be
	// distributed among the participants
	totalLeft := ticketPrice

	// Algorithm sketch: loop starting from lowest amount to highest.
	//
	// If all participants starting at the current one can fill the ticket, then
	// average the contribution of the remaining participants. Otherwise, all
	// participants contribute at least the amount of the current one and we
	// continue trying to split the remaining ticket amount left
	contribs := make([]dcrutil.Amount, nparts)
	for i, max := range remainingMax {
		remainingParts := dcrutil.Amount(nparts - i)
		if remainingParts*max > totalLeft {
			// average the remaining total among the remaining participants
			perPart := dcrutil.Amount(math.Ceil(float64(totalLeft) /
				float64(remainingParts)))
			for j := i; (j < nparts) && (totalLeft > 0); j++ {
				if totalLeft < perPart {
					// due to the division, some participants may contribute
					// slightly less
					contribs[j] += totalLeft
					totalLeft = 0
				} else {
					contribs[j] += perPart
					totalLeft -= perPart
				}
			}
			break
		} else {
			// not enough to split the participation equally, so everyone
			// contributes as much as this participant
			for j := i; j < nparts; j++ {
				contribs[j] += max
				totalLeft -= max
				remainingMax[j] -= max
			}
		}
	}

	// Now account for the pool fees. Split the total pool fees, such that each
	// participant contributes its % to the common pool fee.
	poolFees := make([]dcrutil.Amount, nparts)
	totalPoolFeeLeft := poolFee
	for i := 0; i < nparts; i++ {
		contribPerc := float64(contribs[i]) / float64(ticketPrice)
		fee := dcrutil.Amount(math.Ceil(contribPerc * float64(poolFee)))
		if fee > totalPoolFeeLeft {
			poolFees[i] = totalPoolFeeLeft
			contribs[i] -= totalPoolFeeLeft
			break
		} else {
			poolFees[i] = fee
			totalPoolFeeLeft -= fee
			contribs[i] -= fee
		}
	}

	return contribs, poolFees, nil
}

// SecretNumber is the secret number that each individual  participant chooses.
type SecretNumber []byte

// Hash gives the hash of the secret number, given the hash of a block to use as
// salt. This is usually called with the block id of the mainchain tip.
func (nb SecretNumber) Hash(mainchainHash *chainhash.Hash) SecretNumberHash {
	var res SecretNumberHash
	var calculated []byte

	// note that block hashes are reversed, so the first bytes are the actual
	// random bytes (ending bytes should be a string of 0000s)
	h := blake256.NewSalt(mainchainHash[:16])
	h.Write(nb[:])
	calculated = h.Sum(nil)
	copy(res[:], calculated)

	return res
}

func (nb SecretNumber) String() string {
	return hex.EncodeToString(nb)
}

// Format fulfills the fmt.Formatter interface
func (nb SecretNumber) Format(f fmt.State, c rune) {
	f.Write([]byte(nb.String()))
}

// RandomSecretNumber returns a cryptografically random secret number for use in
// split ticket transactions.
func RandomSecretNumber() SecretNumber {
	res := make(SecretNumber, SecretNbSize)
	_, err := rand.Read(res[:])
	if err != nil {
		// out of entropy. Nothing to do.
		panic(err)
	}

	return res
}

// SecretNumberHash represents the hash of a secret number
type SecretNumberHash [SecretNbHashSize]byte

// Equals checks whether the hashes are equal.
func (h SecretNumberHash) Equals(other SecretNumberHash) bool {
	return bytes.Equal(h[:], other[:])
}

// String converts the secret hash to a string representation.
func (h SecretNumberHash) String() string {
	return hex.EncodeToString(h[:])
}

// LotteryCommitmentHash is the hash that commits all participants of a split
// ticket session to a given lottery result, which decides the voter of the
// session. This commitment is added to the split tx as an OP_RETURN, in order
// for the lottery results to be accountable (ie, a single participant can prove
// whether the posted ticket was the one agreed upon by the rules of the
// lottery).
type LotteryCommitmentHash [LotteryCommitmentHashSize]byte

// CalcLotteryCommitmentHash calculates the lottery commitment hash, given all
// required information.
//
// The lottery commitment hash is calculated as follows:
//   H(secretNbHashes || amounts || voteAddresses)
//
// The mainchainHash[16:] is used as salt to the calculation. The individual
// items are concatenated in slice order.
//
// The number of secrets and amounts **MUST** be the same, otherwise the result
// is undefined
func CalcLotteryCommitmentHash(secretNbHashes []SecretNumberHash,
	amounts []dcrutil.Amount, voteAddresses []dcrutil.Address,
	mainchainHash *chainhash.Hash) *LotteryCommitmentHash {

	var data [8]byte
	var calculated []byte
	res := new(LotteryCommitmentHash)

	// note that block hashes are reversed, so the first bytes are the actual
	// random bytes (ending bytes should be a string of 0000s)
	h := blake256.NewSalt(mainchainHash[:16])

	for _, hash := range secretNbHashes {
		h.Write(hash[:])
	}
	for _, amount := range amounts {
		binary.LittleEndian.PutUint64(data[:], uint64(amount))
		h.Write(data[:])
	}
	for _, addr := range voteAddresses {
		btsAddr := []byte(addr.EncodeAddress())
		h.Write(btsAddr)
	}

	calculated = h.Sum(nil)
	copy(res[:], calculated)

	return res
}

// CalcLotteryResultHash calculates the hash result of the lottery. This value
// is interpreted as a 256 bit uint, such that its value (mod total contribution
// amount) is the coin that should belong to the voter.
func CalcLotteryResultHash(secretNbs []SecretNumber,
	mainchainHash *chainhash.Hash) []byte {

	// hash the secret numbers to get a 256 bit number
	h := blake256.NewSalt(mainchainHash[:16])
	for _, nb := range secretNbs {
		h.Write(nb[:])
	}

	return h.Sum(nil)
}

// CalcLotteryResult discovers the the selected coin and selected index (ie, the
// results of the voter selection lottery) given the input data. Note that
// len(secretNbs) MUST be equal to len(contribAmounts), otherwise the result may
// be undefined.
func CalcLotteryResult(secretNbs []SecretNumber,
	contribAmounts []dcrutil.Amount, mainchainHash *chainhash.Hash) (
	dcrutil.Amount, int) {

	var contribSum dcrutil.Amount
	for _, c := range contribAmounts {
		contribSum += c
	}

	hash := CalcLotteryResultHash(secretNbs, mainchainHash)
	coinBig := big.NewInt(0)
	n := big.NewInt(0)
	n.SetBytes(hash[:])
	m := big.NewInt(int64(contribSum))
	coinBig.Mod(n, m)

	coin := coinBig.Uint64()
	index := 0

	contribSum = 0
	for i, c := range contribAmounts {
		contribSum += c
		if coin < uint64(contribSum) {
			index = i
			break
		}
	}

	return dcrutil.Amount(coin), index
}
