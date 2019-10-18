module github.com/matheusd/dcr-split-ticket-matcher

require (
	github.com/decred/dcrd/blockchain v1.0.2
	github.com/decred/dcrd/blockchain/stake v1.0.2
	github.com/decred/dcrd/certgen v1.0.1
	github.com/decred/dcrd/chaincfg v1.1.1
	github.com/decred/dcrd/chaincfg/chainhash v1.0.1
	github.com/decred/dcrd/crypto/blake256 v1.0.0
	github.com/decred/dcrd/database v1.0.2 // indirect
	github.com/decred/dcrd/dcrec v0.0.0-20180809193022-9536f0c88fa8
	github.com/decred/dcrd/dcrec/secp256k1 v1.0.0
	github.com/decred/dcrd/dcrjson v1.0.0
	github.com/decred/dcrd/dcrutil v1.1.1
	github.com/decred/dcrd/gcs v1.0.2 // indirect
	github.com/decred/dcrd/hdkeychain v1.1.0
	github.com/decred/dcrd/rpcclient v1.0.1
	github.com/decred/dcrd/txscript v1.0.1
	github.com/decred/dcrd/wire v1.1.0
	github.com/decred/dcrdata v2.1.3+incompatible
	github.com/decred/dcrwallet/rpc/walletrpc v0.1.0
	github.com/decred/dcrwallet/wallet v1.0.0
	github.com/decred/slog v1.0.0
	github.com/go-ini/ini v1.48.0
	github.com/golang/protobuf v1.1.0
	github.com/gorilla/websocket v1.2.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/mattn/go-gtk v0.0.0-20190930150717-0423bc8d46fb
	github.com/mattn/go-pointer v0.0.0-20190911064623-a0a44394634f // indirect
	github.com/pkg/errors v0.8.1
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.14.0
	gopkg.in/ini.v1 v1.48.0 // indirect
)

go 1.12
