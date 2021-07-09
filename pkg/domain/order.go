package domain

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/kyokan/smallbridge/pkg/pb"
	"github.com/kyokan/smallbridge/pkg/conv"
)




const(         
	posted = iota  //自增值
	intentfill
	done
	invaild
)

type Order struct {
	ID             string
	UserPubkey     *btcec.PublicKey
	MakerAsset      AssetID
	TakerAsset      AssetID
	Quantity       	float32
	Price           float32
	CancelAt       	string
	Memo           	string
	CreatedAt      	string
	Status         	int32
}



type FillIntent struct {
	ID             string
	UserPubkey     string
	CpPubKey 	   string
	OrderId        string
	Quantity       float32
	Price          float32
	CancelAt       string
	Memo           string
	CreatedAt      string
	MakerAsset     string
	TakerAsset     string
	Status         int
}



type FillConfirmInfo struct {
	UserPubKey string
  	CpPubKey   string
    OrderId    string
    IntentId   string
  	HTLCHash   string
    MakerAddress string
    TakerAddress string
    Memo       string
  	CreatedAt  string
  	Signature  string
  	Settled    bool
  	PHash      string
  	Timeout    string
  	Quantity    float32
  	MakerChain string
    TakerChain string
    AssetId    string
    ChainId    string
}

func (o *Order) Protobuffer() (*pb.OrderV1) {
	return &pb.OrderV1{
		OrderId:        o.ID,
		UserPubKey:     conv.PubkeyToStr(o.UserPubkey),
		MakerAsset:     string(o.MakerAsset),
		TakerAsset:		string(o.TakerAsset),
		Quantity:       o.Quantity,
		Price:          o.Price,
		CancelAt:       o.CancelAt,
		Memo:           o.Memo,
		CreatedAt:      o.CreatedAt,
		Signature:      nil,
	}
}

func OrderFromPB(pbo *pb.OrderV1) (*Order, error) {
	pubkey, err := conv.PubkeyFromStr(pbo.UserPubKey)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:             pbo.OrderId,
		UserPubkey:     pubkey,
		MakerAsset:     AssetID(pbo.MakerAsset),
		TakerAsset: 	AssetID(pbo.TakerAsset),
		Quantity:       pbo.Quantity,
		Price:          pbo.Price,
		CancelAt:       pbo.CancelAt,
		Memo:           pbo.Memo,
		CreatedAt:      pbo.CreatedAt,
		Status:		 	invaild,	
	}, nil
}



func FillIntentFromPB(pbf *pb.NotifyFillIntentRequestV1) (*FillIntent, error) {
	return &FillIntent{
		ID:             pbf.IntentId,
		UserPubkey:     pbf.UserPubKey,
		CpPubKey:     	pbf.CpPubKey,
		OrderId: 		pbf.OrderId,
		Quantity:       pbf.Quantity,
		CancelAt:       pbf.CancelAt,
		Memo:           pbf.Memo,
		CreatedAt:      pbf.CreatedAt,
		MakerAsset:     pbf.MakerAsset,
		TakerAsset:     pbf.TakerAsset,
		Price:          pbf.Price,
	}, nil
}


