package clearing

import (
	"math/big"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kyokan/smallbridge/pkg/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"context"
	"time"
	"github.com/pkg/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/kyokan/smallbridge/pkg"
	"github.com/btcsuite/btcd/rpcclient"
)

var wei, _ = new(big.Float).SetString("10000000000000000")

type ETHHTLCManager struct {
	address      string
	rpcURL       string
	chainID      string
	htlcContract *eth.HTLC
	privKey      *btcec.PrivateKey
	ethClient    *ethclient.Client
}

func NewETHHTLCManager(privKey *btcec.PrivateKey, cfg *pkg.ETHConfig) *ETHHTLCManager {
	return &ETHHTLCManager{
		address: cfg.ContractAddress,
		rpcURL:  cfg.RPCUrl,
		chainID: cfg.ChainID,
		privKey: privKey,
	}
}

func (e *ETHHTLCManager) Start() error {
	conn, err := ethclient.Dial(e.rpcURL)
	if err != nil {
		return err
	}
	e.ethClient = conn
	htlc, err := eth.NewHTLC(common.HexToAddress(e.address), conn)
	if err != nil {
		return err
	}
	e.htlcContract = htlc
	logger.Info("started ETH htlc manager")
	return nil
}

func (e *ETHHTLCManager) Stop() error {
	e.ethClient.Close()
	return nil
}

func (e *ETHHTLCManager) ChainID() string {
	return e.chainID
}

func (e *ETHHTLCManager) Broadcast(hash [32]byte, receiver *btcec.PublicKey, amount *big.Float) ([]byte, error) {

	logger.Info("broadcasting transaction", "chain_id", e.chainID)
	auth := bind.NewKeyedTransactor(e.privKey.ToECDSA())
	auth.Value, _ = amount.Mul(amount, wei).Int(nil)
	auth.GasLimit = 2000000
	auth.GasPrice = big.NewInt(2000000000)
	senderAddress := crypto.PubkeyToAddress(*e.privKey.PubKey().ToECDSA())
	receiverAddress := crypto.PubkeyToAddress(*receiver.ToECDSA())
	blk, err := e.ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	height := blk.Number()
	logger.Info("sending ETH HTLC", "gas-limit", auth.GasLimit, "receiver", receiverAddress, "hash", hexutil.Encode(hash[:]), "value", auth.Value.Text(10))
	tx, err := e.htlcContract.NewContract(auth, receiverAddress, hash, big.NewInt(time.Now().Add(24 * time.Hour).Unix()))
	if err != nil {
//	fmt.Printf("\n1ssssssss3333356\n")
		return nil, err
	}
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timed out")
		case <-tick.C:
			logger.Info("polling for logs", "hash", tx.Hash())
			iter, err := e.htlcContract.FilterLogHTLCNew(&bind.FilterOpts{
				Start:   height.Uint64(),
				Context: context.Background(),
			}, nil, []common.Address{senderAddress}, nil)

			if err != nil {
				return nil, err
			}

			if iter.Next() {
				logger.Info("broadcast new HTLC event", "id", hexutil.Encode(iter.Event.ContractId[:]), "hash", hexutil.Encode(hash[:]))
				return iter.Event.ContractId[:], nil
			}
		}
	}
}

func (*ETHHTLCManager) Claim(preimage []byte) error {
	panic("implement me")
}

func (e *ETHHTLCManager) Redeemed(id []byte) (bool, [32]byte, error) {
	logger.Info("checking redemption status", "chain_id", e.chainID, "tx_hash", hexutil.Encode(id))
	var idArr [32]byte
	copy(idArr[:], id)
	contract, err := e.htlcContract.GetContract(nil, idArr)
	var emptyPreimage [32]byte
	if err != nil {
		return false, emptyPreimage, err
	}

	if contract.Preimage != emptyPreimage {
		return true, contract.Preimage, nil
	}

	return false, emptyPreimage, nil
}

func (e *ETHHTLCManager) Redeemable(id []byte) (bool, error) {
	logger.Info("checking redeemability", "chain_id", e.chainID, "tx_hash", hexutil.Encode(id))
	var idArr [32]byte
	copy(idArr[:], id)
	contract, err := e.htlcContract.GetContract(nil, idArr)
	if err != nil {
		return false, err
	}
	// simply check for existence
	var emptyAddress common.Address
	return contract.Sender != emptyAddress, nil
}

func (e *ETHHTLCManager) Redeem(id []byte, preimage [32]byte) (error) {
	logger.Info("redeeming HTLC", "tx_hash", hexutil.Encode(id), "chain_id", e.chainID, "preimage", hexutil.Encode(preimage[:]))
	var idArr [32]byte
	copy(idArr[:], id)
	auth := bind.NewKeyedTransactor(e.privKey.ToECDSA())
	auth.GasLimit = 4000000
	_, err := e.htlcContract.Withdraw(auth, idArr, preimage)
	return err
}


func (e *ETHHTLCManager) GetAddress() string {
	return crypto.PubkeyToAddress(*e.privKey.PubKey().ToECDSA()).String()
}


func (e *ETHHTLCManager) GetBtcClient() *rpcclient.Client {
	return nil
}

func (e *ETHHTLCManager) GetEthClient() *ethclient.Client {
	return e.ethClient
}