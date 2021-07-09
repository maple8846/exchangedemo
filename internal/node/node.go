package node

import (
	"github.com/kyokan/smallbridge/pkg"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"github.com/kyokan/smallbridge/pkg/pb"
	log "github.com/inconshreveable/log15"
	"github.com/kyokan/smallbridge/internal/storage"
	"context"
	"github.com/kyokan/smallbridge/internal/p2p"
	"github.com/kyokan/smallbridge/internal/clearing"
	//"github.com/kyokan/smallbridge/pkg/btc"
)

var logger = log.New("module", "node")

type Node struct {
	addr      string
	port      int
	dbUrl     string
	store     storage.Storage
	cancel    context.CancelFunc
	book      *p2p.PeerBook
	ourPeerId *pkg.PeerID
	settler   *clearing.Settler

	pkg.BaseService
}

func NewNode(ctx context.Context, book *p2p.PeerBook, ourPeerId *pkg.PeerID, store storage.Storage, settler *clearing.Settler, addr string, port int) *Node {
	n := &Node{
		addr:      addr,
		port:      port,
		store:     store,
		ourPeerId: ourPeerId,
		book:      book,
		settler:   settler,
	}
	n.Ctx, n.cancel = context.WithCancel(ctx)
	return n
}

func (n *Node) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.addr, n.port))
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pb.RegisterNodeServer(srv, &nodeServer{
		store:     n.store,
		ourPeerId: n.ourPeerId,
		book:      n.book,
		Settler:   n.settler,
	})
	go srv.Serve(lis)

	go func() {
		<-n.Ctx.Done()
		srv.Stop()
		return
	}()

	logger.Info("started gRPC endpoint", "addr", n.addr, "port", n.port)
	return nil
}

func (n *Node) Stop() error {
	n.cancel()
	return nil
}
