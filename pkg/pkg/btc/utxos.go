package btc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"encoding/hex"
	)

func FindBestUTXOs(utxos []btcjson.ListUnspentResult, amt btcutil.Amount) ([]btcjson.ListUnspentResult, error) {
	// first, try to find a single UTXO that's larger than the amount

	for _, utxo := range utxos {
		cmp, err := btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return nil, err
		}
		if amt < cmp {
			return []btcjson.ListUnspentResult{utxo}, nil
		}
	}

	// next try to find a set of UTXOs that sum to the amount
	panic("not implemented yet")
}

func DecodePkScript(script string) ([]byte, error) {
	src := []byte(script)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}