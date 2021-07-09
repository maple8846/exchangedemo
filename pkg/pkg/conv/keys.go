package conv

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/btcsuite/btcutil/base58"
	"fmt"
	"strings"
	"github.com/pkg/errors"
)

const SmallbridgeKey = "sbrdg"

func PubkeyToStr(pub *btcec.PublicKey) string {
	return fmt.Sprintf("%s%s", SmallbridgeKey, base58.Encode(pub.SerializeCompressed()))
}

func PubkeyFromStr(str string) (*btcec.PublicKey, error) {
	if !strings.HasPrefix(str, SmallbridgeKey) {
		return nil, errors.New("invalid prefix")
	}
	str = str[len(SmallbridgeKey):]
	return btcec.ParsePubKey(base58.Decode(str), btcec.S256())
}

func PrivkeyFromHex(hex string) (*btcec.PrivateKey, error) {
	b, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), b)
	return priv, nil
}
