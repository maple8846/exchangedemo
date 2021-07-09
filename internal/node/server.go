package node

import (
	"github.com/kyokan/smallbridge/internal/storage"
	"github.com/kyokan/smallbridge/pkg/pb"
	"golang.org/x/net/context"
	"github.com/kyokan/smallbridge/pkg"
	"github.com/kyokan/smallbridge/internal/p2p"
	"github.com/kyokan/smallbridge/internal/clearing"
	"github.com/kyokan/smallbridge/pkg/domain"
	"github.com/pkg/errors"
	"github.com/kyokan/smallbridge/pkg/conv"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"fmt"
	//"strconv"
		"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

const(         
	posted = iota  //自增值
	intentfill
	done
	cancel
)

const(         
	intentposted = iota  //自增值
	intentdone
)

type nodeServer struct {
	store     storage.Storage
	Settler   *clearing.Settler
	ourPeerId *pkg.PeerID
	book      *p2p.PeerBook
}




//const MaxSupportAsset int = 3

func (n *nodeServer) Identify(ctx context.Context, req *pb.IdentifyRequestV1) (*pb.IdentifyResponseV1, error) {
	logger.Info("received identity message", "peer_id", req.PeerId)

	theirPeerId, err := pkg.PeerIDFromStr(req.PeerId)
	if err != nil {
		return nil, err
	}
	go n.book.AddPeer(theirPeerId)

	return &pb.IdentifyResponseV1{
		PeerId: n.ourPeerId.String(),
	}, nil
}

func (n *nodeServer) NotifyOrder(ctx context.Context, req *pb.NotifyOrderRequestV1) (*pb.NotifyOrderResponseV1, error) {

	logger.Info("received order notification", "pubkey", req.Order.UserPubKey)
	order, err := domain.OrderFromPB(req.Order)
	if err != nil {
		
		return nil, err
	}

	_, err = n.store.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	//n.settler.NotifyOrder(order.ID)
	return &pb.NotifyOrderResponseV1{}, nil
}

func (n *nodeServer) NotifyFill(ctx context.Context, req *pb.NotifyFillRequestV1) (*pb.NotifyFillResponseV1, error) {
	logger.Info("received fill notification", "order_id", req.OrderId)
	if len(req.PreimageHash) != 32 {
		return nil, errors.New("mal-formed preimage")
	}
	pub, err := conv.PubkeyFromStr(req.UserPubkey)
	if err != nil {
		return nil, err
	}

	var h [32]byte
	copy(h[:], req.PreimageHash)
	n.Settler.NotifyFill(clearing.FillInfo{
		OrderID:      req.OrderId,
		PreimageHash: h,
		TxHash:       req.TxHash,
		Counterparty: pub,
		IntentId:     req.IntentId,
	})
	return &pb.NotifyFillResponseV1{}, nil
}

func (n *nodeServer) NotifySettle(ctx context.Context, req *pb.NotifySettleRequestV1) (*pb.NotifySettleResponseV1, error) {
	logger.Info("received settlement notification", "order_id", req.OrderId)
	n.Settler.Finalize(req.OrderId, req.TxHash)
	return &pb.NotifySettleResponseV1{}, nil
}



func (n *nodeServer) NotifyFillIntent(ctx context.Context, req *pb.NotifyFillIntentRequestV1) (*pb.NotifyFillIntentResponseV1, error) {
	//logger.Info("received order notification", "pubkey", req.Order.UserPubKey)
	fillintent, _ := domain.FillIntentFromPB(req)
	err := n.store.UpdateOrderStatus(req.OrderId,intentfill)
	if err != nil {
		return nil, err
	}
	

	err = n.store.RecordFillIntentOrder(fillintent)
	if err != nil {
		return nil, err
	}

	return &pb.NotifyFillIntentResponseV1{}, nil
}



func (n *nodeServer) NotifyFillConfirm(ctx context.Context, req *pb.NotifyFillConfirmRequestV1) (*pb.NotifyFillConfirmResponseV1, error) {
	
	err := n.store.UpdateOrderStatus(req.OrderId,done)
    
	if err != nil {
		return nil, err
	}
	  
	err = n.store.UpdateIntentStatus(req.IntentId, intentdone)	
	if err != nil { 
		return nil, err
	}

	n.Settler.NotifyOrder(req.OrderId, req.IntentId)
	return &pb.NotifyFillConfirmResponseV1{}, nil
}


func (n *nodeServer) NotifyDeleteIntent(ctx context.Context, req *pb.NotifyDeleteIntentRequestV1) (*pb.NotifyDeleteIntentResponseV1, error) {
	
	err := n.store.UpdateIntentStatus(req.IntentId,intentdone)	
	
	if err !=nil {
		return nil,err
	}
	return &pb.NotifyDeleteIntentResponseV1{},nil
}

func (n *nodeServer) NotifyDeleteOrder(ctx context.Context, req *pb.NotifyDeleteOrderRequestV1) (*pb.NotifyDeleteOrderResponseV1, error) {
	
	err := n.store.UpdateOrderStatus(req.OrderId,cancel)	
	
	if err !=nil {
		return nil,err
	}
	return &pb.NotifyDeleteOrderResponseV1{},nil
}


func (n *nodeServer) NotifyGetBalances(ctx context.Context, req *pb.NotifyGetBalancesRequestV1) (*pb.NotifyGetBalancesResponseV1, error) {
	assetId := req.AssetId
	chainId := req.ChainId
	var singleresponse *pb.NotifyGetBalancesResponse
	var finalresponse  *pb.NotifyGetBalancesResponseV1
	finalresponse = new(pb.NotifyGetBalancesResponseV1)
	var assetName []string
	var chainName []string
	quantity := 2

	zero := big.NewInt(0)
	switch assetId {
		case "BTC-BTC-TESTNET":

			mgr:=(n.Settler).Sw.Mgrs[chainId]
			addressstr := mgr.GetAddress()

			address,_:=btcutil.DecodeAddress(addressstr, &chaincfg.TestNet3Params)
			client := mgr.GetBtcClient()

			account,_:= client.GetAccount(address)
			amount, _ := client.GetBalance(account)
			fmt.Printf("\nXXXXXXX\n",amount.String())
			
			//if amount != 0 {
				singleresponse = new(pb.NotifyGetBalancesResponse)
				singleresponse.Quantity = string(quantity)
				singleresponse.AssetId  = assetId
				singleresponse.ChainId  = chainId
				singleresponse.Address  = address.EncodeAddress()
				singleresponse.AddressQuantity = amount.String()  
				finalresponse.Responses = append(finalresponse.Responses,singleresponse)
			//} 
		case "ETH-ETH-RINKEBY":

			mgr:=(n.Settler).Sw.Mgrs[chainId]
			address := mgr.GetAddress()
			//fmt.Printf("eeeeeeeeeeeeeee",address)
			client := mgr.GetEthClient()
			amount,err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
    		if err != nil {
    			return nil,err
    		}   

    		if amount.Cmp(zero) != 0 {
				
				singleresponse = new(pb.NotifyGetBalancesResponse)
				singleresponse.Quantity = string(quantity)
				singleresponse.AssetId  = assetId
				singleresponse.ChainId  = chainId
				singleresponse.Address  = address
				singleresponse.AddressQuantity = amount.String()
				finalresponse.Responses = append(finalresponse.Responses,singleresponse)
			
			} 
		
		default:
			quantity, _:= n.store.GetAssetId(&assetName,&chainName)
			for i:=0;i<len(assetName);i++{
					switch assetName[i] {

							case "BTC-BTC-TESTNET":
								mgr:=(n.Settler).Sw.Mgrs[chainName[i]]
								address := mgr.GetAddress()
								client := mgr.GetBtcClient()
								amount, _ := client.GetBalance(address)
								
								if amount != 0 {
										
										singleresponse = new(pb.NotifyGetBalancesResponse)
										singleresponse.Quantity = string(quantity)
										singleresponse.AssetId  = assetName[i]
										singleresponse.ChainId  = chainName[i]
										singleresponse.Address  = address
										singleresponse.AddressQuantity = string(amount)
										finalresponse.Responses = append(finalresponse.Responses,singleresponse)
								
								}

							case "ETH-ETH-RINKEBY":
									mgr:=(n.Settler).Sw.Mgrs[chainId]
									address := mgr.GetAddress()
									client := mgr.GetEthClient() 
									amount,err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
    									if err != nil {
    										//fmt.Printf("zzzzzzzzzz\n")
    										return nil,err
    								}   

					    		if amount.Cmp(zero) != 0 {
					
									singleresponse = new(pb.NotifyGetBalancesResponse)
									singleresponse.Quantity = string(quantity)
									singleresponse.AssetId  = assetName[i]
									singleresponse.ChainId  = chainName[i]
									singleresponse.Address  = address
									singleresponse.AddressQuantity = amount.String()
									finalresponse.Responses = append(finalresponse.Responses,singleresponse)
				
								} 
						}
			
			} 
		}
	return finalresponse,nil
}



func (n *nodeServer) NotifyFillConfirmLocal(ctx context.Context, req *pb.NotifyFillConfirmRequestV1) (*pb.NotifyFillConfirmResponseV1, error) {
	
	err := n.store.UpdateOrderStatus(req.OrderId,done)

	if err != nil {
		return nil, err
	}
	   	
//	n.Settler.NotifyOrder(req.OrderId, req.IntentId)
	return &pb.NotifyFillConfirmResponseV1{}, nil
}

