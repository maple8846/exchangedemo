package clearing

import (
	"context"
	"github.com/kyokan/smallbridge/internal/p2p"
	"github.com/kyokan/smallbridge/internal/storage"
	"github.com/kyokan/smallbridge/pkg/pb"
	"github.com/kyokan/smallbridge/pkg/domain"
	"github.com/btcsuite/btcd/btcec"
	"github.com/kyokan/smallbridge/pkg/conv"
	"time"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/kyokan/smallbridge/pkg"
	"github.com/kyokan/smallbridge/pkg/util"
	"crypto/sha256"
	"math/big"
	"fmt"
)

type tradeSide int
type watchType int

const (
	poster tradeSide = iota
	responder

	redeemed watchType = iota
	redeemable
)

type htlcWatch struct {
	OrderID      string
	ourChainID   string
	theirChainID string
	OurTxHash    []byte
	theirTxHash  []byte
}

type PendingRedemption struct {
	OrderID  string
	chainID  string
	Preimage [32]byte
}

type htlcPost struct {
	OrderID string
	txHash  []byte
}

type FillInfo struct {
	OrderID      string
	Counterparty *btcec.PublicKey
	PreimageHash [32]byte
	TxHash       []byte
}


type FillIntent struct {
	OrderID      string
	Counterparty *btcec.PublicKey
}



type Settler struct {
	orderCh chan string
	intentId string
	fillCh  chan FillInfo
	watched chan htlcWatch

	pending chan PendingRedemption
	posts   chan htlcPost

	cancel       context.CancelFunc
	book         *p2p.PeerBook
	storage      storage.Storage
	Sw           *HTLCSwitch
	ourPubkey    *btcec.PublicKey
	WatchedHTLCs []htlcWatch

	PendingRedemptions map[string]PendingRedemption

	pkg.BaseService
}


func NewSettler(book *p2p.PeerBook, storage storage.Storage, sw *HTLCSwitch, ourPubkey *btcec.PublicKey) *Settler {
	res := &Settler{
		orderCh:            make(chan string),
		fillCh:             make(chan FillInfo),
		watched:            make(chan htlcWatch),
		pending:            make(chan PendingRedemption),
		posts:              make(chan htlcPost),
		book:               book,
		storage:            storage,
		Sw:                 sw,
		ourPubkey:          ourPubkey,
		PendingRedemptions: make(map[string]PendingRedemption),
		intentId:          "",
	}
	res.Ctx, res.cancel = context.WithCancel(context.Background())
	return res
}

func (s *Settler) Start() error {
	go s.fillLoop()
	go s.htlcWatcher()
	return nil
}

func (s *Settler) Stop() error {
	s.cancel()
	return nil
}

func (s *Settler) fillLoop() {
	for {
		select {
		case id := <-s.orderCh:
			go s.fillOrder(id)
		case info := <-s.fillCh:
			go s.settleOrder(info)
		case <-s.Ctx.Done():
			return
		}
	}
}

func (s *Settler) htlcWatcher() {
	tick := time.NewTicker(10 * time.Second)

	for {
		select {
		case w := <-s.watched:
			s.WatchedHTLCs = append(s.WatchedHTLCs, w)
		case r := <-s.pending:
			s.PendingRedemptions[r.OrderID] = r
		case <-tick.C:
			s.checkHTLCs()
		case post := <-s.posts:
			s.redeemHTLC(post)
		case <-s.Ctx.Done():
			return
		}
	}
}

func (s *Settler) NotifyOrder(orderID string,intentID string) {
	s.intentId = intentID
	s.orderCh <- orderID
	
}



func (s *Settler) NotifyFill(info FillInfo) {
	s.fillCh <- info
}

func (s *Settler) Finalize(orderID string, txHash []byte) {
	s.posts <- htlcPost{
		OrderID: orderID,
		txHash:  txHash,

	}
}

func (s *Settler) fillOrder(id string) {
	logger.Info("filling order", "id", id)

	order, _, mgr, err := s.fetchOrderAndManager(id, responder)
	if err != nil {
		return
	}

	if conv.PubkeyToStr(order.UserPubkey) == conv.PubkeyToStr(s.ourPubkey){

		return
	}
    preimage := util.Rand32()
	hash := sha256.Sum256(preimage[:])
	total := order.Quantity*order.Price

	txHash, err := mgr.Broadcast(hash, order.UserPubkey, big.NewFloat(float64(total)))
	if err != nil {
		logger.Error("failed to broadcast transaction", "order_id", order.ID, "err", err)
		return
	}

	waitingChainID, err := s.chainIDForSide(order, poster)
	if err != nil {
		logger.Error("failed to fetch other side", "err", err)
		return
	}

	s.pending <- PendingRedemption{
		OrderID:  order.ID,
		chainID:  waitingChainID,
		Preimage: preimage,
	}

	peer := s.book.PeerByPubkey(order.UserPubkey)
	_, err = peer.Client.NotifyFill(context.Background(), &pb.NotifyFillRequestV1{
		OrderId:      order.ID,
		PreimageHash: hash[:],
		TxHash:       txHash,
		UserPubkey:   conv.PubkeyToStr(s.ourPubkey),
	})
	if err != nil {
		logger.Error("received error broadcasting fill notification", "err", err)
	}
}

func (s *Settler) settleOrder(info FillInfo) {
	fmt.Printf("开始结算")
	logger.Info("settling order", "id", info.OrderID)
	order, chainID, mgr, err := s.fetchOrderAndManager(info.OrderID, poster)
	if err != nil {
		return
	}


	txHash, err := mgr.Broadcast(info.PreimageHash, info.Counterparty, big.NewFloat(float64(order.Quantity)))
	if err != nil {
		logger.Error("failed to broadcast transaction", "order_id", order.ID, "err", err)
		return
	}

	theirChainID, err := s.chainIDForSide(order, responder)
	if err != nil {
		logger.Error("failed to fetch other side", "err", err)
		return
	}

	s.watched <- htlcWatch{
		ourChainID:   chainID,
		theirChainID: theirChainID,
		OurTxHash:    txHash,
		theirTxHash:  info.TxHash,
	}

	peer := s.book.PeerByPubkey(info.Counterparty)
	_, err = peer.Client.NotifySettle(context.Background(), &pb.NotifySettleRequestV1{
		OrderId: order.ID,
		TxHash:  txHash,
	})
	if err != nil {
		logger.Error("received error broadcasting settlement notification", "err", err)
		return
	}
	

	htlcHashinfo := &storage.HTLCHashInfo{
			HTLCHash:     hexutil.Encode(info.TxHash[:]),      
			PHash:     	  hexutil.Encode(info.PreimageHash[:]), 
			Timeout:        "0",
	}

	cpchtlcHashinfo := &storage.HTLCHashInfo{
			HTLCHash:     hexutil.Encode(info.TxHash[:]),      
			PHash:     	  hexutil.Encode(info.PreimageHash[:]), 
			Timeout:        "0",
	}
	s.storage.RecordFillInfo(htlcHashinfo,cpchtlcHashinfo,order.ID,s.intentId)
}

func (s *Settler) checkHTLCs() {
	if len(s.WatchedHTLCs) == 0 {
		logger.Info("no HTLCs to watch")
	}

	watched := make([]htlcWatch, 0)

	for _, htlc := range s.WatchedHTLCs {
		ourMgr := s.Sw.Manager(htlc.ourChainID)
		if ourMgr == nil {
			// should never happen, implies bug
			panic("htlc manager does not exist")
		}

		redeemed, Preimage, err := ourMgr.Redeemed(htlc.OurTxHash)
		if err != nil {
			logger.Error("failed to fetch redemption status", "err", err)
			watched = append(watched, htlc)
			continue
		}

		if !redeemed {
			logger.Info("htlc not redeemed, trying again later", "tx_hash", hexutil.Encode(htlc.OurTxHash))
			watched = append(watched, htlc)
			continue
		}

		theirMgr := s.Sw.Manager(htlc.theirChainID)
		if ourMgr == nil {
			// should never happen, implies bug
			panic("htlc manager does not exist")
		}

		logger.Info("got HTLC redemption", "preimage", hexutil.Encode(Preimage[:]), "chain_id", htlc.ourChainID, "tx_hash", hexutil.Encode(htlc.OurTxHash))
		err = theirMgr.Redeem(htlc.theirTxHash, Preimage)
		if err != nil {
			logger.Error("failed to redeem HTLC", "err", err)
		}
		logger.Info("redeemed HTLC", "chain_id", htlc.theirChainID)
	}

	s.WatchedHTLCs = watched
}

func (s *Settler) redeemHTLC(posted htlcPost) {
	pending := s.PendingRedemptions[posted.OrderID]
	mgr := s.Sw.Manager(pending.chainID)
	if err := mgr.Redeem(posted.txHash, pending.Preimage); err != nil {
		logger.Error("failed to redeem htlc", "err", err)
		return
	}
	logger.Info("successfully redeemed HTLC", "chain_id", pending.chainID)
	delete(s.PendingRedemptions, posted.OrderID)
}

func (s *Settler) fetchOrderAndManager(id string, side tradeSide) (*domain.Order, string, HTLCManager, error) {
	order, err := s.storage.FindOrder(id)
	if order == nil {
		return nil, "", nil, err
	}
	if err != nil {
		logger.Warn("failed to fetch order", "err", err)
		return nil, "", nil, err
	}
	chainID, err := s.chainIDForSide(order, side)
	if err != nil {
		logger.Warn("failed to fetch order chain ID", "order_id", order.ID, "err", err)
		return nil, "", nil, err
	}
	mgr := s.Sw.Manager(chainID)
	if mgr == nil {
		logger.Info("no registered HTLC manager found", "chain_id", chainID, "order_id", order.ID)
		return nil, "", nil, err
	}

	return order, chainID, mgr, nil
}

func (s *Settler) chainIDForSide(order *domain.Order, side tradeSide) (string, error) {
	var assetId domain.AssetID

	switch side {
	case poster:
		assetId = order.MakerAsset
	case responder:
		assetId = order.TakerAsset
	}

	chainID, err := s.storage.ChainIDForAsset(assetId)
	if err != nil {
		return "", err
	}

	return chainID, nil
}


