package btc

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
		"crypto/sha256"
)

// Generates a BIP-199 HTLC script.
// OP_IF
//     OP_SHA256 <hash> OP_EQUALVERIFY OP_DUP OP_HASH160 <instant pubkey hash>
// OP_ELSE
//     <num> OP_CSV OP_DROP OP_DUP OP_HASH160 <delayed pubkey hash>
// OP_ENDIF
// OP_EQUALVERIFY
// OP_CHECKSIG
func GenHTLCScript(hash [32]byte, instantPub *btcec.PublicKey, delayedPub *btcec.PublicKey) ([]byte, error) {
	instantHash := Hash160(instantPub.SerializeCompressed())
	delayedHash := Hash160(delayedPub.SerializeCompressed())

	bldr := txscript.NewScriptBuilder()
	bldr.AddOp(txscript.OP_IF)
	bldr.AddOp(txscript.OP_SHA256)
	bldr.AddData(hash[:])
	bldr.AddOp(txscript.OP_EQUALVERIFY)
	bldr.AddOp(txscript.OP_DUP)
	bldr.AddOp(txscript.OP_HASH160)
	bldr.AddData(instantHash)
	bldr.AddOp(txscript.OP_ELSE)
	bldr.AddOp(txscript.OP_7)
	bldr.AddOp(txscript.OP_CHECKSEQUENCEVERIFY)
	bldr.AddOp(txscript.OP_DROP)
	bldr.AddOp(txscript.OP_DUP)
	bldr.AddOp(txscript.OP_HASH160)
	bldr.AddData(delayedHash)
	bldr.AddOp(txscript.OP_ENDIF)
	bldr.AddOp(txscript.OP_EQUALVERIFY)
	bldr.AddOp(txscript.OP_CHECKSIG)

	return bldr.Script()
}

func GenHTLCRedemption(preimage [32]byte) ([]byte, error) {
	bldr := txscript.NewScriptBuilder()
	bldr.AddData(preimage[:])
	bldr.AddOp(txscript.OP_TRUE)
	return bldr.Script()
}

// OP_TRUE preimage pubkey
// preimage pubkey
// hash hash pubkey
// pubkey pubkey
// ha sh ash pubkey

func Hash160(in []byte) []byte {
	sha := sha256.New()
	sha.Write(in)
	rmd := ripemd160.New()
	rmd.Write(sha.Sum(nil))
	return rmd.Sum(nil)
}