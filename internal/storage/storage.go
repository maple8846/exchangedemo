package storage

import (
		"github.com/kyokan/smallbridge/pkg"
	"github.com/kyokan/smallbridge/pkg/domain"
"github.com/kyokan/smallbridge/pkg/pb"
)

type Storage interface {
	pkg.Service
	CreateOrder(order *domain.Order) (*domain.Order, error)
	FindOrder(id string) (*domain.Order, error)
	ChainIDForAsset(assetID domain.AssetID) (string, error)
	FindOrderbyUser(req *pb.GetOrderRequestV1)(*pb.GetOrderResponseV1, error)
	UpdateOrderStatus(id string, status int) (error)
	GetFillIntentOrder(req *pb.GetIntentRequestV1) (*pb.GetIntentResponseV1,error)
	RecordFillIntentOrder(fillIntent *domain.FillIntent) (error) 
	GetAssetId(idarray *[]string,chainid *[]string) (int,error)
	RecordFillInfo(para *HTLCHashInfo,cpppara *HTLCHashInfo, orderid string, intentid string) (error)
	GetFillInfo(id string) (*pb.GetFillResponseV1,error)
	GetFillconfirmInfo(id string) (*pb.FillConfirmResponseV1,error)
	UpdateIntentStatus(id string, status int) (error)
	RecordLocalFillInfo(para *HTLCHashInfo, orderid string, pubkey string) (error)
	RecordCppFillInfo(para *HTLCHashInfo, intentid string) (error)
	MakerRecordLocalFillInfo(para *HTLCHashInfo, intentId string)(error)
}
