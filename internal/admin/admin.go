package admin

import (
	"github.com/kyokan/smallbridge/pkg"
	"context"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"github.com/kyokan/smallbridge/pkg/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"net/http"
	log "github.com/inconshreveable/log15"
	"github.com/kyokan/smallbridge/internal/p2p"
	"github.com/kyokan/smallbridge/internal/storage"
	"github.com/btcsuite/btcd/btcec"
)

var logger = log.New("module", "admin")

type Admin struct {
	rpcAddr  string
	rpcPort  int
	httpAddr string
	httpPort int
	cancel   context.CancelFunc
	book     *p2p.PeerBook
	storage  storage.Storage
	pubkey   *btcec.PublicKey
	localnodeId *pkg.PeerID

	pkg.BaseService
}

func NewAdmin(ctx context.Context, book *p2p.PeerBook, storage storage.Storage, pubkey *btcec.PublicKey, rpcAddr string, rpcPort int, httpAddr string, httpPort int, localnodeId *pkg.PeerID) *Admin {
	a := &Admin{
		rpcAddr:  rpcAddr,
		rpcPort:  rpcPort,
		httpAddr: httpAddr,
		httpPort: httpPort,
		book:     book,
		storage:  storage,
		pubkey:   pubkey,
		localnodeId: localnodeId,
	}
	a.Ctx, a.cancel = context.WithCancel(ctx)
	return a
}

func (a *Admin) Start() error {
	if err := a.startGRPC(); err != nil {
		return err
	}

	if err := a.startProxy(); err != nil {
		return err
	}

	return nil
}

func (a *Admin) Stop() error {
	a.cancel()
	return nil
}

func (a *Admin) startGRPC() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.rpcAddr, a.rpcPort))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	pb.RegisterAdminServer(srv, &adminServer{
		Book:    a.book,
		storage: a.storage,
		pubkey:  a.pubkey,
		localnodeId: a.localnodeId,
	})
	go func() {
		if err := srv.Serve(lis); err != nil {
			logger.Error("gRPC server error", "err", err)
		}
	}()

	go func() {
		<-a.Ctx.Done()
		srv.Stop()
		logger.Info("gRPC endpoint shut down")
	}()

	logger.Info("started gRPC endpoint", "addr", a.rpcAddr, "port", a.rpcPort)
	return nil
}

func (a *Admin) startProxy() error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	rpcEndpoint := fmt.Sprintf("%s:%d", a.rpcAddr, a.rpcPort)
	if err := pb.RegisterAdminHandlerFromEndpoint(a.Ctx, mux, rpcEndpoint, opts); err != nil {
		return err
	}
	httpEndpoint := fmt.Sprintf("%s:%d", a.httpAddr, a.httpPort)
	srv := &http.Server{Addr: httpEndpoint, Handler: mux}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error("reverse proxy error", "err", err)
		}
	}()

	go func() {
		<-a.Ctx.Done()
		srv.Shutdown(a.Ctx)
		logger.Info("reverse proxy shut down")
	}()

	logger.Info("started reverse proxy", "addr", a.httpAddr, "port", a.httpPort)
	return nil
}
