package admin

import (
	"github.com/kyokan/smallbridge/pkg/pb"
	"golang.org/x/net/context"
	"github.com/kyokan/smallbridge/internal/p2p"
	"github.com/kyokan/smallbridge/pkg"
	"github.com/kyokan/smallbridge/internal/storage"
	"github.com/kyokan/smallbridge/pkg/domain"
	"github.com/kyokan/smallbridge/pkg/conv"
	"github.com/btcsuite/btcd/btcec"
	"github.com/beevik/guid"
	"time"
	"strconv"
	"github.com/ethereum/go-ethereum/common/hexutil"
		"github.com/ethereum/go-ethereum/crypto"
	"github.com/kyokan/smallbridge/pkg/btc"
)

const( 
	posted = iota  //自增值
	intentfill
	done
	cancel
)

const( 
	intentpost = iota  //自增值
	intentdone
)
type adminServer struct {
	Book    *p2p.PeerBook
	storage storage.Storage
	pubkey  *btcec.PublicKey
	localnodeId *pkg.PeerID
}



func (a *adminServer) ConnectPeer(ctx context.Context, req *pb.ConnectPeerRequestV1) (*pb.ConnectPeerResponseV1, error) {
	peerId, err := pkg.PeerIDFromStr(req.PeerId)
	if err != nil {
		return nil, err
	}
	id, err := a.Book.AddPeer(peerId)
	if err != nil {
		return nil, err
	}

	return &pb.ConnectPeerResponseV1{
		PeerId: id.String(),
	}, nil
}


func (a *adminServer) PostOrder(ctx context.Context, req *pb.PostOrderRequestV1) (*pb.PostOrderResponseV1, error) {
	

	t := time.Now()
	tempcreated := uint64(t.Unix())
	created := strconv.FormatUint(tempcreated, 10)
	order := &domain.Order{
		UserPubkey:     a.pubkey,
		MakerAsset:   	domain.AssetID(req.MakerAsset),
		TakerAsset: 	domain.AssetID(req.TakerAsset),
		Quantity:       req.Quantity,
		Price:          req.Price,
		Memo:           req.Memo,
		CancelAt:       req.CancelAt,
		CreatedAt:		created,	
		Status:         posted,	
	}		

	order, err := a.storage.CreateOrder(order)

	if err != nil {
		return nil, err
	}


	go func() {  
		notifyOrderV1 := &pb.OrderV1{
				OrderId: 		order.ID,
				UserPubKey:		conv.PubkeyToStr(order.UserPubkey),
				MakerAsset:		string(order.MakerAsset),
				TakerAsset:		string(order.TakerAsset),
				Quantity:		order.Quantity,
				Price:			order.Price,
				CancelAt:		order.CancelAt,
				CreatedAt:		created,
				Memo:			order.Memo,
				Signature:			nil,
		}
		errs := a.Book.Broadcast(func(p *p2p.Peer) error {
			notifyReq := &pb.NotifyOrderRequestV1{
			 		Order: notifyOrderV1,
			}
			_, err := p.Client.NotifyOrder(context.Background(), notifyReq)
			if err != nil {
				logger.Error("caught error broadcasting to peer", "peer_id", p.ID, "err", err)
			}

			return err
		})

		if len(errs) > 0 {
			logger.Warn("received errors during broadcast", "count", len(errs))
		}
	}()
    

			
	return &pb.PostOrderResponseV1{	
		OrderId:        order.ID,
		CreatedAt: order.CreatedAt,
		UserPubkey: hexutil.Encode(order.UserPubkey.SerializeCompressed()),
		MakerAsset: string(order.MakerAsset),
		TakerAsset: string(order.TakerAsset),
		Quantity:   order.Quantity,
		Price: 		order.Price,
		CancelAt:  	order.CancelAt,
		Memo:       order.Memo,
		Signature:  nil,
	}, nil
}


func (a *adminServer) GetOrder(ctx context.Context, req *pb.GetOrderRequestV1) (*pb.GetOrderResponseV1,error) { 
    
    var getOrderResponseV1  *pb.GetOrderResponseV1  
	
	getOrderResponseV1, err := a.storage.FindOrderbyUser(req)
	if err != nil {
		return nil, err
	}
	return  getOrderResponseV1,nil
}



func (a *adminServer) GetNodeInfo(ctx context.Context, req *pb.GetNodeInfoRequestV1) (*pb.GetNodeInfoResponseV1,error) {    
    response := &pb.GetNodeInfoResponseV1{}  	
	response.UserPubKey = a.Book.OurPeerId.String()
	return  response,nil
}





func (a *adminServer) FillIntent(ctx context.Context, req *pb.FillIntentRequestV1) (*pb.FillIntentResponseV1,error) { 
	
	err := a.storage.UpdateOrderStatus(req.OrderId, intentfill)

    order, _ := a.storage.FindOrder(req.OrderId)

	if err != nil {
		return nil, err
	}
	guid := guid.New()
	t := time.Now()
	now := t.Unix()
	created := strconv.FormatUint(uint64(now), 10)
	fillIntent:= &domain.FillIntent{
		ID:				guid.String(),
		UserPubkey:     conv.PubkeyToStr(a.pubkey),
		CpPubKey: 	    conv.PubkeyToStr(order.UserPubkey),
		OrderId:        req.OrderId,	
		Quantity:       order.Quantity,
		CancelAt:       req.CancelAt,
		Memo:           req.Memo,
		CreatedAt:      created,
		MakerAsset:     string(order.MakerAsset),
		TakerAsset:     string(order.TakerAsset),
		Status:			intentpost,
		Price:          order.Price,
	}
    err = a.storage.RecordFillIntentOrder(fillIntent)
	if err != nil {
		return nil, err
	}	
		go func() {
		errs := a.Book.Broadcast(func(p *p2p.Peer) error {
			notifyReq := &pb.NotifyFillIntentRequestV1{
				IntentId:			fillIntent.ID,
				UserPubKey:			fillIntent.UserPubkey,
				CpPubKey:			fillIntent.CpPubKey,
				OrderId: 			fillIntent.OrderId,
				Quantity:			fillIntent.Quantity,
				CreatedAt:			fillIntent.CreatedAt,
				CancelAt:			fillIntent.CancelAt,
				Memo:				fillIntent.Memo,
				Price:              fillIntent.Price,
				Signature:			nil,
				MakerAsset:         fillIntent.MakerAsset,
				TakerAsset:         fillIntent.TakerAsset,
				Status:             strconv.Itoa(fillIntent.Status),
			}

			_, err := p.Client.NotifyFillIntent(context.Background(), notifyReq)

			if err != nil {
				logger.Error("caught error broadcasting to peer", "peer_id", p.ID, "err", err)
			}
			return err
		})

		if len(errs) > 0 {
			logger.Warn("received errors during broadcast", "count", len(errs))
		}
	}()



    userpubkey,_ := conv.PubkeyFromStr(fillIntent.UserPubkey)	
	cppubkey,_ := conv.PubkeyFromStr(fillIntent.CpPubKey)

	return &pb.FillIntentResponseV1{
		IntentId:        fillIntent.ID,
		UserPubKey:		 hexutil.Encode(userpubkey.SerializeCompressed()),
		CpPubKey:		 hexutil.Encode(cppubkey.SerializeCompressed()),
		OrderId:		 fillIntent.OrderId,
		Quantity:	     fillIntent.Quantity,
		CreatedAt:	     fillIntent.CreatedAt,
		CancelAt: 		 fillIntent.CancelAt,
		Memo:			 fillIntent.Memo,
		Signature:       nil,
		MakerAsset:      fillIntent.MakerAsset,
		TakerAsset:      fillIntent.TakerAsset,
		Status:          strconv.Itoa(fillIntent.Status),
		Price:         	 fillIntent.Price,
	}, nil
}




func (a *adminServer) FillConfirm(ctx context.Context, req *pb.FillConfirmRequestV1) (*pb.FillConfirmResponseV1, error) {
	var tempres *pb.FillConfirmResponseV1
	err := a.storage.UpdateOrderStatus(req.OrderId, done)	
	if err != nil { 
		return nil, err
	}

	err =a.storage.UpdateIntentStatus(req.IntentId, intentdone)	
	if err != nil { 
		return nil, err
	}
	
	go func() {
		errs := a.Book.Broadcast(func(p *p2p.Peer) error {
			notifyReq := &pb.NotifyFillConfirmRequestV1{
				IntentId:	req.IntentId,  				
				OrderId:	req.OrderId, 
			}
			_, err := p.Client.NotifyFillConfirm(context.Background(), notifyReq)

			if err != nil {
				logger.Error("caught error broadcasting to peer", "peer_id", p.ID, "err", err)
			}
			return err
		})

		if len(errs) > 0 {
			logger.Warn("received errors during broadcast", "count", len(errs))
		}
	}()	
    
    time.Sleep(time.Duration(15)*time.Second)
	tempres, err = a.storage.GetFillconfirmInfo(req.OrderId)
	if tempres == nil {
		return nil, err
	}	
	switch tempres.AssetId {
		case "BTC-BTC-TESTNET":
            makeraddress,_ := conv.PubkeyFromStr(tempres.UserPubKey)
            takeraddress,_ := conv.PubkeyFromStr(tempres.CpPubKey)
			tempres.MakerAddress = crypto.PubkeyToAddress(*makeraddress.ToECDSA()).String()
			tempres.TakerAddress = crypto.PubkeyToAddress(*takeraddress.ToECDSA()).String()
		case "ETH-ETH-RINKEBY":
			tempmakeradress,_ :=conv.PubkeyFromStr(tempres.UserPubKey)
			temptakeradress,_ :=conv.PubkeyFromStr(tempres.CpPubKey)
			makeraddress,_:= btc.PubkeyToAddress(tempmakeradress)
			takeraddress,_:= btc.PubkeyToAddress(temptakeradress)
			tempres.MakerAddress = makeraddress.EncodeAddress()
			tempres.TakerAddress = takeraddress.EncodeAddress()
	}

	return tempres,nil
}


func (a *adminServer)GetFillIntent(ctx context.Context, req *pb.GetIntentRequestV1)(*pb.GetIntentResponseV1, error){
	
	var getFillResponseV1  *pb.GetIntentResponseV1  	
	getFillResponseV1, err := a.storage.GetFillIntentOrder(req)	
	
	if err !=nil {
		return nil,err
	}
	
	return getFillResponseV1, nil
}

func (a *adminServer)CancelIntent(ctx context.Context, req *pb.CancelIntentRequestV1)(*pb.CancelIntentResponseV1, error){

	err := a.storage.UpdateIntentStatus(req.IntentId, intentdone)	
	
	if err !=nil {
		return nil,err
	}
	go func() {
		errs := a.Book.Broadcast(func(p *p2p.Peer) error {
			notifyReq := &pb.NotifyDeleteIntentRequestV1{
				IntentId:	req.IntentId,  				   
			}
			_, err := p.Client.NotifyDeleteIntent(context.Background(), notifyReq)

			if err != nil {
				logger.Error("caught error broadcasting to peer", "peer_id", p.ID, "err", err)
			}
			return err
		})

		if len(errs) > 0 {
			logger.Warn("received errors during broadcast", "count", len(errs))
		}
	}()	
	
	return &pb.CancelIntentResponseV1{}, nil
}



func (a *adminServer)CancelOrder(ctx context.Context, req *pb.CancelOrderRequestV1)(*pb.CancelOrderResponseV1, error){
	err := a.storage.UpdateOrderStatus(req.OrderId, cancel)		
	if err !=nil {
		return nil,err
	}
	go func() {
		errs := a.Book.Broadcast(func(p *p2p.Peer) error {
			notifyReq := &pb.NotifyDeleteOrderRequestV1{
				OrderId:	req.OrderId,  				   
			}
			_, err := p.Client.NotifyDeleteOrder(context.Background(), notifyReq)

			if err != nil {
				logger.Error("caught error broadcasting to peer", "peer_id", p.ID, "err", err)
			}
			return err
		})

		if len(errs) > 0 {
			logger.Warn("received errors during broadcast", "count", len(errs))
		}
	}()	
	
	return &pb.CancelOrderResponseV1{}, nil
}


func (a *adminServer)GetBalances(ctx context.Context, req *pb.GetBalancesRequestV1) (*pb.GetBalancesResponseV1, error){
	peerInfo := &pb.ConnectPeerRequestV1{}
	peerInfo.PeerId = a.localnodeId.String()	
	_,err:=a.ConnectPeer(ctx,peerInfo)
	
	if err != nil {
		return nil,err
	}
	
	notifyReq := &pb.NotifyGetBalancesRequestV1{
			ChainId: 	req.ChainId,
    		AssetId:	req.AssetId,		   
		}	
	
	response,err := a.Book.Clients[a.localnodeId.String()].Client.NotifyGetBalances(context.Background(), notifyReq)		
	var finalresponse *pb.GetBalancesResponseV1
	var singleresponse *pb.GetBalancesResponse
	finalresponse = new(pb.GetBalancesResponseV1) 
	length := len(response.Responses)
		
	for i:=0;i<length;i++{
		singleresponse = new(pb.GetBalancesResponse)
		singleresponse.ChainId = response.Responses[i].ChainId
		singleresponse.AssetId = response.Responses[i].AssetId
		singleresponse.Quantity = response.Responses[i].Quantity
		singleresponse.Address = response.Responses[i].Address
		singleresponse.AddressQuantity = response.Responses[i].AddressQuantity
		finalresponse.Responses = append(finalresponse.Responses,singleresponse)
	}
	return finalresponse, err
}


func (a *adminServer)GetFill(ctx context.Context, req *pb.GetFillRequestV1)(*pb.GetFillResponseV1, error){
	
	response, err := a.storage.GetFillInfo(req.OrderId)
	if err!= nil{
		return nil,err
	}
	return response,nil
}