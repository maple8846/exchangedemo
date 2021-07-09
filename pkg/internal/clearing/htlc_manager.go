package clearing

import (
	"math/big"
	"github.com/btcsuite/btcd/btcec"
	log "github.com/inconshreveable/log15"
	"github.com/kyokan/smallbridge/pkg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
)

var logger = log.New("module", "clearing")

type HTLCManager interface {
	pkg.Service

	ChainID() string
	Broadcast(hash [32]byte, receiver *btcec.PublicKey, amount *big.Float) ([]byte, error)
	Claim(preimage []byte) error
	Redeemed(txHash []byte) (bool, [32]byte, error)
	Redeemable(txHash []byte) (bool, error)
	Redeem(txHash []byte, preimage [32]byte) (error)
	GetAddress() string
	GetBtcClient() *rpcclient.Client
	GetEthClient() *ethclient.Client
}
