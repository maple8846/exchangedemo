package container

import (
	"github.com/kyokan/smallbridge/pkg"
	"context"
	"github.com/kyokan/smallbridge/internal/storage"
	"github.com/kyokan/smallbridge/internal/node"
	"github.com/kyokan/smallbridge/internal/admin"
	"os"
	"github.com/kyokan/smallbridge/pkg/util"
	log "github.com/inconshreveable/log15"
	"github.com/kyokan/smallbridge/internal/p2p"
	"github.com/kyokan/smallbridge/pkg/conv"
	"github.com/kyokan/smallbridge/internal/clearing"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kyokan/smallbridge/pkg/btc"
)


func Start(cfg pkg.Config) error {
	priv, err := conv.PrivkeyFromHex(cfg.PrivateKeyHex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	store, err := storage.NewPGStorage(cfg.DBUrl)
	if err != nil {
		return err
	}
	if err := store.Start(); err != nil {
		return err
	}

	ourPeerId := &pkg.PeerID{
		Pubkey: priv.PubKey(),
		Addr:   cfg.Node.RPCAddr,
		Port:   cfg.Node.RPCPort,
	}
    //由于结算过程是在node处进行，getbalance操作需要在node处处理，而node模块未开放http接口，需要admin模块调用node rpc进行处理，因此将本地id保存在全局变量中

	book := p2p.NewPeerBook(ourPeerId)

	sw := clearing.NewHTLCSwitch()
	for _, eCfg := range cfg.ETHConfig {
		mgr := clearing.NewETHHTLCManager(priv, &eCfg)
		if err := sw.RegisterManager(mgr); err != nil {
			return err
		}
	}
	for _, bCfg := range cfg.BTCConfig {
		mgr := clearing.NewBTCHTLCManager(priv.PubKey(), &bCfg, priv)
		if err := sw.RegisterManager(mgr); err != nil {
			return err
		}
	}
	if err := sw.Start(); err != nil {
		return err
	}

	settler := clearing.NewSettler(book, store, sw, priv.PubKey())
	if err := settler.Start(); err != nil {
		return err
	}

	n := node.NewNode(ctx, book, ourPeerId, store, settler, cfg.Node.RPCAddr, cfg.Node.RPCPort)
	if err := n.Start(); err != nil {
		return err
	}
	a := admin.NewAdmin(ctx, book, store, priv.PubKey(), cfg.Admin.RPCAddr, cfg.Admin.RPCPort, cfg.Admin.HTTPAddr, cfg.Admin.HTTPPort, ourPeerId)
	if err := a.Start(); err != nil {
		return err
	}

	btcaddr, err := btc.PubkeyToAddress(priv.PubKey())
	if err != nil {
		return err
	}
	ethAddr := crypto.PubkeyToAddress(*priv.PubKey().ToECDSA())

	log.Info("started", "peer_id", ourPeerId, "btc_address", btcaddr.EncodeAddress(), "eth_address", ethAddr)
	util.AwaitTermination(func() {
		log.Info("interrupted, shutting down")
		cancel()
		os.Exit(0)
	})

	return nil
}
