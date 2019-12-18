module github.com/matheusd/dcr-split-ticket-matcher

require (
	github.com/dchest/blake256 v1.0.0
	github.com/decred/dcrd/blockchain v1.0.2
	github.com/decred/dcrd/blockchain/stake v1.0.2
	github.com/decred/dcrd/certgen v1.0.1
	github.com/decred/dcrd/chaincfg v1.1.1
	github.com/decred/dcrd/chaincfg/chainhash v1.0.1
	github.com/decred/dcrd/database v1.0.2 // indirect
	github.com/decred/dcrd/dcrec v0.0.0-20180816212643-20eda7ec9229
	github.com/decred/dcrd/dcrec/edwards v0.0.0-20180816212643-20eda7ec9229 // indirect
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
	github.com/go-ini/ini v1.30.0
	github.com/golang/protobuf v1.1.0
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/gorilla/websocket v1.2.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jtolds/gls v4.2.1+incompatible // indirect
	github.com/mattn/go-gtk v0.0.0-20191030024613-af2e013261f5
	github.com/mattn/go-pointer v0.0.0-20190911064623-a0a44394634f // indirect
	github.com/pkg/errors v0.8.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/smartystreets/goconvey v0.0.0-20180222194500-ef6db91d284a // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/net v0.0.0-20180808004115-f9ce57c11b24
	google.golang.org/grpc v1.14.0
)

replace github.com/btcsuite/goleveldb => github.com/matheusd/goleveldb v0.0.0-20190109102927-255f9b4fc43c

go 1.12
