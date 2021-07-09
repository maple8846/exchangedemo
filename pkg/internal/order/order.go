package orderbook

import (
	"strconv"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)


addressChecksumLen：=4
func ValidateAddress(address string) bool {
    pubKeyHash := Base58Decode([]byte(address))
    actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
    version := pubKeyHash[0]
    pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
    targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
    return bytes.Compare(actualChecksum, targetChecksum) == 0
}


func NewBlockchain(nodeID string) *Blockchain {
    dbFile := fmt.Sprintf(dbFile, nodeID)
    if dbExists(dbFile) == false {
        fmt.Println("No existing blockchain found. Create one first.")
        os.Exit(1)
    }
    var tip []byte
    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil {
        log.Panic(err)
    }
    err = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        tip = b.Get([]byte("l"))
        return nil
    })
    if err != nil {
        log.Panic(err)
    }
    bc := Blockchain{tip, db}
    return &bc
}


func getBalance(address, nodeID string) {
    if !ValidateAddress(address) {
        log.Panic("ERROR: Address is not valid")
    }
    bc := NewBlockchain(nodeID)
    UTXOSet := UTXOSet{bc}
    defer bc.db.Close()
    balance := 0
    pubKeyHash := Base58Decode([]byte(address))
    pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
    UTXOs := UTXOSet.FindUTXO(pubKeyHash)
    for _, out := range UTXOs {
        balance += out.Value
    }
    return balance
}
type Blockchain struct {
    tip []byte
    Db  *bolt.DB
}
type UTXOSet struct {
    Blockchain *Blockchain
}



func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
    var UTXOs []TXOutput
    db := u.Blockchain.db
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            outs := DeserializeOutputs(v)
            for _, out := range outs.Outputs {
                if out.IsLockedWithKey(pubKeyHash) {   //如果公钥hash能够解锁UTXO集合中的交易输出，则该交易内部的值属于该地址
                    UTXOs = append(UTXOs, out)   //添加到UTXOS中
                }
            }
        }
        return nil
    })
    if err != nil {
        log.Panic(err)
    }
    return UTXOs
}
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
    return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}




func QuryLiquidityProof(str string, amount big.Int, address string) bool{
	switch str {
	case string(BTC):
		balance = getBalance(address,MainBtc)
		if balance >= amount {
			return true
		}
		else {
			return false
		}
		return BTC, nil
	case string(ETH):
		conn, err := ethclient.Dial("https://mainnet.infura.io")
    	if err != nil {
    		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    	}
    	err = client.Call(&result, "eth_getBalance", account[0], "latest")
    	account := common.HexToAddress(address)
		balance, err := client.BalanceAt(context.Background(), account, nil)
		if balance >= amount {
			return true
		}
		else {
			return false
		}
	case string(USDT):
		return USDT, nil
	
}


