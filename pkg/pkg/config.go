package pkg

import (
	"io/ioutil"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Config struct {
	DBUrl         string      `mapstructure:"db_url"`
	Admin         AdminConfig `mapstructure:"admin"`
	Node          NodeConfig  `mapstructure:"node"`
	ETHConfig     []ETHConfig `mapstructure:"eth"`
	BTCConfig     []BTCConfig `mapstructure:"btc"`
	PrivateKeyHex string      `mapstructure:"private_key"`
}

type AdminConfig struct {
	HTTPAddr string `mapstructure:"http_addr"`
	HTTPPort int    `mapstructure:"http_port"`
	RPCAddr  string `mapstructure:"rpc_addr"`
	RPCPort  int    `mapstructure:"rpc_port"`
}

type NodeConfig struct {
	RPCAddr string `mapstructure:"rpc_addr"`
	RPCPort int    `mapstructure:"rpc_port"`
}

type ETHConfig struct {
	ContractAddress string `mapstructure:"contract_address"`
	ChainID         string `mapstructure:"chain_id"`
	RPCUrl          string `mapstructure:"rpc_url"`
}

type BTCConfig struct {
	ChainID     string `mapstructure:"chain_id"`
	RPCUrl      string `mapstructure:"rpc_url"`
	RPCCertFile string `mapstructure:"rpc_cert_file"`
	RPCUsername string `mapstructure:"rpc_username"`
	RPCPassword string `mapstructure:"rpc_password"`
}

const DefaultConfig = `db_url = "postgres://localhost:5432/smallbridge?sslmode=disable"
private_key="%s"

[admin]
http_addr="0.0.0.0"
http_port=8098
rpc_addr="127.0.0.1"
rpc_port=8081

[node]
rpc_addr="127.0.0.1"
rpc_port=8082

[[eth]]
contract_address=""
chain_id="ETH-RINKEBY"
rpc_url="http://127.0.0.1:8545"

[[btc]]
chain_id="BTC-TESTNET"
rpc_url="127.0.0.1:18334"
rpc_cert_file="./dev/rpc.cert"
rpc_username="user"
rpc_password="password"`

func WriteDefaultConfigFile(name string) error {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return err
	}

	hex := hexutil.Encode(privKey.Serialize())

	return ioutil.WriteFile(name, []byte(fmt.Sprintf(DefaultConfig, hex)), 0744)
}
