package pkg

const (
	// ProtocolVersion records protocol-incompatible changes to the matching
	// engine. Buyers and the matching server must be running the same protocol
	// version in order for the matching process to proceed.
	// History of Versions:
	// v1: Initial version
	// v2: Modified algo that calculates the lottery commiment hash to include
	// the amounts and voter addresses
	// v3: Modified pool fee to be proportional by contribution amount
	ProtocolVersion = 3
)

// These are the individual version numbers
const (
	MajorVersion = 0
	MinorVersion = 6
	PatchVersion = 2
)

// Version is the package version. Remember to modify the previous constants as
// well on version bumps.
const Version = "0.6.2+pre"
