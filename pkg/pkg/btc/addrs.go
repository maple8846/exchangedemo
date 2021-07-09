package btc

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func PubkeyToAddress(pub *btcec.PublicKey) (*btcutil.AddressPubKey, error) {
	btcaddr, err := btcutil.NewAddressPubKey(pub.SerializeUncompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}
	return btcaddr, nil
}