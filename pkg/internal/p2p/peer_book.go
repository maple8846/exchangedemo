package p2p

import (
	"github.com/kyokan/smallbridge/pkg/pb"
	"google.golang.org/grpc"
	"context"
	"github.com/kyokan/smallbridge/pkg"
	"sync"
	"time"
	"github.com/kyokan/smallbridge/pkg/util"
	log "github.com/inconshreveable/log15"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"fmt"
)

var logger = log.New("module", "peer-book")


type BroadcastFunc func(p *Peer) error

type Peer struct {
	conn   *grpc.ClientConn
	Client pb.NodeClient
	ID     *pkg.PeerID
}

type PeerBook struct {
	Clients   map[string]*Peer
	byPubkey  map[string]string
	mtx       sync.RWMutex
	OurPeerId *pkg.PeerID
}

func NewPeerBook(ourPeerId *pkg.PeerID) *PeerBook {

	return &PeerBook{
		Clients:   make(map[string]*Peer),
		byPubkey:  make(map[string]string),
		OurPeerId: ourPeerId,
	}
}

func (p *PeerBook) AddPeer(peerId *pkg.PeerID) (*pkg.PeerID, error) {
	logger.Info("adding new peer", "id", peerId)
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if _, ok := p.Clients[peerId.String()]; ok {
		logger.Info("already have peer", "peer_id", peerId)
		return peerId, nil
	}

	conn, err := grpc.Dial(util.CombineHostPort(peerId.Addr, peerId.Port), grpc.WithInsecure())
	if err != nil {
		logger.Error("failed to dial peer", "peer_id", peerId, "err", err)
		return nil, err
	}

	client := pb.NewNodeClient(conn)
	deadline := time.Now().Add(5 * time.Second)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	res, err := client.Identify(ctx, &pb.IdentifyRequestV1{
		PeerId:          p.OurPeerId.String(),
		ProtocolVersion: pkg.ProtocolVersion,
	})
	fmt.Printf("\nour peer id xxxxxxxx", res.PeerId)
	if err != nil {
		logger.Error("peer failed to identify itself", "peer_id", peerId, "err", err)
		return nil, err
	}

	id, err := pkg.PeerIDFromStr(res.PeerId)
	if err != nil {
		logger.Error("peer sent mal-formed peer ID", "sent_id", res.PeerId, "err", err)
		return nil, err
	}

	peer := &Peer{
		conn:   conn,
		Client: client,
		ID:     id,
	}

	p.byPubkey[hexutil.Encode(id.Pubkey.SerializeCompressed())] = id.String()
	p.Clients[id.String()] = peer
	return id, nil
}

func (p *PeerBook) PeerByID(id *pkg.PeerID) (*Peer) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	return p.Clients[id.String()]
}

func (p *PeerBook) PeerByPubkey(pubkey *btcec.PublicKey) (*Peer) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	id := p.byPubkey[hexutil.Encode(pubkey.SerializeCompressed())]
	return p.Clients[id]
}

func (p *PeerBook) Broadcast(f BroadcastFunc) []error {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
    
	/*广播的时候不要发给自己,记录本地Id*/
	localId := p.OurPeerId
	/**/
	errors := make([]error, 0)
	for _, p := range p.Clients {
		/*广播的时候不要发给自己*/
		if p.ID == localId {
			continue
		}
		/*广播的时候不要发给自己*/
		err := f(p)

		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}



