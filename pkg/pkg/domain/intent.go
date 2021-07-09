//Should check with Matt if this is helpful.
package domain

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/kyokan/smallbridge/pkg/conv"
	"github.com/kyokan/smallbridge/pkg/pb"
)

//Intent refers to intents to fill standing orders.
type Intent struct {
	IntentId   string
	UserPubkey *btcec.PublicKey
	CpPubkey   *btcec.PublicKey
	OrderId    OrderId
	Quantity   *big.Float
	CancelAt   uint64
	CreatedAt  uint64
	Memo       string
	Signature  string //Should this by some other type?
}

//Protobuffer reformats the intent.
func (in *Intent) Protobuffer() *pb.IntentV1 {
	return &pb.IntentV1{
		IntentId:   in.IntentId,
		UserPubkey: conv.PubkeyToStr(in.UserPubkey),
		CpPubkey:   conv.PubkeyToStr(in.CpPubkey),
		OrderId:    in.OrderId,
		Quantity:   in.Quantity.String(),
		CancelAt:   in.CancelAt,
		CreatedAt:  in.CreatedAt,
		Memo:       in.Memo,
		Signature:  nil,
	}
}

func IntentFromPB(pbo *pb.IntentV1) (*Intent, error) {
	pubkey, err := conv.PubkeyFromStr(pbo.UserPubkey)
	if err != nil {
		return nil, err
	}

	return &Intent{
		IntentId:   pbo.IntentId,
		UserPubkey: UserPubKey,
		CpPubkey:   CpPubKey,
		OrderId:    pbo.OrderId,
		Quantity:   conv.BigFloatFromStr(pbo.Quantity),
		CancelAt:   pbo.CancelAt,
		CreatedAt:  pbo.CreatedAt,
		Memo:       pbo.Memo,
	}, nil
}
