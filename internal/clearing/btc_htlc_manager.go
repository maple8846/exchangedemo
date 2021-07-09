package clearing

import (
	"math/big"
	"github.com/btcsuite/btcd/btcec"
	"io/ioutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
	"github.com/kyokan/smallbridge/pkg/btc"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"bytes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/kyokan/smallbridge/pkg"
	//"github.com/kyokan/smallbridge/pkg/conv"
	"github.com/ethereum/go-ethereum/ethclient"
)

const fee = 10000

type BTCHTLCManager struct {
	rpcUrl      string
	rpcCertFile string
	rpcUsername string
	rpcPassword string
	chainID     string
	OurPubkey   *btcec.PublicKey
	privKey     *btcec.PrivateKey

	client *rpcclient.Client
}

func NewBTCHTLCManager(ourPubkey *btcec.PublicKey, cfg *pkg.BTCConfig, privKey *btcec.PrivateKey) *BTCHTLCManager {
	return &BTCHTLCManager{
		rpcUrl:      cfg.RPCUrl,
		rpcCertFile: cfg.RPCCertFile,
		rpcUsername: cfg.RPCUsername,
		rpcPassword: cfg.RPCPassword,
		chainID:     cfg.ChainID,
		OurPubkey:   ourPubkey,
		privKey:     privKey,
	}
}

func (m *BTCHTLCManager) Start() error {
	cert, err := ioutil.ReadFile(m.rpcCertFile)
	if err != nil {
		return err
	}

	cfg := &rpcclient.ConnConfig{
		Host:         m.rpcUrl,
		Endpoint:     "ws",
		User:         m.rpcUsername,
		Pass:         m.rpcPassword,
		Certificates: cert,
	}
	client, err := rpcclient.New(cfg, &rpcclient.NotificationHandlers{})
	if err != nil {
		logger.Error("failed to start RPC client", "err", err, "url", m.rpcUrl)
		return err
	}
	m.client = client
	logger.Info("started BTC htlc manager")
	return nil
}

func (m *BTCHTLCManager) Stop() error {
	m.client.Shutdown()
	m.client.WaitForShutdown()
	return nil
}

func (m *BTCHTLCManager) ChainID() string {
	return m.chainID
}

func (m *BTCHTLCManager) Broadcast(hash [32]byte, receiver *btcec.PublicKey, amount *big.Float) ([]byte, error) {
	unspent, err := m.client.ListUnspent()
	if err != nil {
		return nil, err
	}
	if len(unspent) == 0 {
		return nil, errors.New("no outputs to spend")
	}

	amtFloat, _ := amount.Float64()
	amt, err := btcutil.NewAmount(amtFloat)
	ourAddr, err := btc.PubkeyToAddress(m.OurPubkey)
	tx := wire.NewMsgTx(wire.TxVersion)
	usable, err := btc.FindBestUTXOs(unspent, amt)
	if err != nil {
		return nil, err
	}

	total := btcutil.Amount(0)
	currInputs := make([]*wire.TxIn, 0)
	for _, utxo := range usable {
		amt, err := btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return nil, err
		}
		total += amt
		h, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		outpoint := wire.NewOutPoint(h, utxo.Vout)
		in := wire.NewTxIn(outpoint, nil, nil)
		currInputs = append(currInputs, in)
		tx.AddTxIn(in)
	}

	changeScript, err := txscript.PayToAddrScript(ourAddr)
	if err != nil {
		return nil, err
	}

	changeAmount := total - amt - fee
	change := wire.NewTxOut(int64(changeAmount), changeScript)
	tx.AddTxOut(change)

	destScript, err := btc.GenHTLCScript(hash, receiver, m.OurPubkey)
	if err != nil {
		return nil, err
	}
	out := wire.NewTxOut(int64(amt), destScript)
	logger.Info("here is my value", "val", out.Value)
	tx.AddTxOut(out)

	tx, ok, err := m.client.SignRawTransaction(tx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("signing failed")
	}

	var b bytes.Buffer
	tx.Serialize(&b)

	txHash, err := m.client.SendRawTransaction(tx, true)
	logger.Info("broadcast tx", "hash", hexutil.Encode(txHash.CloneBytes()))
	if err != nil {
		return nil, err
	}
	return txHash.CloneBytes(), nil	
}

func (*BTCHTLCManager) Claim(preimage []byte) error {
	panic("implement me")
}

func (m *BTCHTLCManager) Redeemed(txHash []byte) (bool, [32]byte, error) {
	logger.Info("checking redemption status", "chain_id", m.chainID, "tx_hash", hexutil.Encode(txHash))
	var emptyPreimage [32]byte
	h, err := chainhash.NewHash(txHash)
	if err != nil {
		return false, emptyPreimage, err
	}
	tx, err := m.client.GetTransaction(h)
	if err != nil {
		return false, emptyPreimage, err
	}
	out, err := m.client.GetTxOut(h, tx.Details[0].Vout, false)
	if out == nil {
		return false, emptyPreimage, nil
	}

	return true, emptyPreimage, nil
}

func (m *BTCHTLCManager) Redeemable(txHash []byte) (bool, error) {
	logger.Info("checking redeemability", "chain_id", m.chainID, "tx_hash", hexutil.Encode(txHash))
	h, err := chainhash.NewHash(txHash)
	if err != nil {
		return false, err
	}
	tx, err := m.client.GetTransaction(h)
	if err != nil {
		return false, err
	}
	out, err := m.client.GetTxOut(h, tx.Details[0].Vout, false)
	return out != nil, err
}

func (m *BTCHTLCManager) Redeem(txHash []byte, preimage [32]byte) ([]byte,error) {
	logger.Info("redeeming", "chain_id", m.chainID, "tx_hash", hexutil.Encode(txHash))
	h, err := chainhash.NewHash(txHash)
	if err != nil {
		return nil,err
	}
	inTx, err := m.client.GetRawTransaction(h)
	if err != nil {
		return nil,err
	}
	ourAddr, err := btc.PubkeyToAddress(m.OurPubkey)
	if err != nil {
		return nil,err
	}
	if err != nil {
		return nil,err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	point := wire.NewOutPoint(h, 1)
	in := wire.NewTxIn(point, nil, nil)
	tx.AddTxIn(in)
	changeScript, err := txscript.PayToAddrScript(ourAddr)
	if err != nil {
		return nil,err
	}
	htlc := inTx.MsgTx().TxOut[1]
	changeAmount := htlc.Value - fee
	change := wire.NewTxOut(int64(changeAmount), changeScript)
	tx.AddTxOut(change)
	redemptionScript, err := btc.GenHTLCRedemption(preimage)
	if err != nil {
		return nil,err
	}
	sigScript, err := txscript.SignatureScript(tx, 0, htlc.PkScript, txscript.SigHashAll, m.privKey, true)
	if err != nil {
		return nil,err
	}

	var buf bytes.Buffer
	buf.Write(sigScript)
	buf.Write(redemptionScript)
	tx.TxIn[0].SignatureScript = buf.Bytes()
	tHash, err := m.client.SendRawTransaction(tx, false)
	logger.Info("redeemed", hexutil.Encode(tHash.CloneBytes()))
	return tHash.CloneBytes() ,err
}

func (m *BTCHTLCManager) GetAddress() string {
	addr,_:= btc.PubkeyToAddress(m.OurPubkey)
	return addr.EncodeAddress()
}



func (m *BTCHTLCManager) GetBtcClient() *rpcclient.Client {
	return m.client
}

func (m *BTCHTLCManager) GetEthClient() *ethclient.Client {
	return nil
}