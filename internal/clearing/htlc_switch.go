package clearing

import (
	"sync"
	"github.com/pkg/errors"
)

type HTLCSwitch struct {
	Mgrs    map[string]HTLCManager
	mtx     sync.RWMutex
	started bool
}

func (s *HTLCSwitch) Start() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, mgr := range s.Mgrs {
		if err := mgr.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (s *HTLCSwitch) Stop() error {
	var failed bool

	for cid, mgr := range s.Mgrs {
		if err := mgr.Stop(); err != nil {
			failed = true
		    logger.Error("failed to stop manager", "chain_id", cid, "err", err)
		}
	}

	if failed {
		return errors.New("some HTLC managers failed to stop; check logs")
	}

	return nil
}

func NewHTLCSwitch() *HTLCSwitch {
	return &HTLCSwitch{
		Mgrs: make(map[string]HTLCManager),
	}
}

func (s *HTLCSwitch) RegisterManager(mgr HTLCManager) error {
	if s.started {
		// don't want to dynamically start newly-added managers
		return errors.New("cannot add new managers after start")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	cid := mgr.ChainID()
	if cid == "" {
		return errors.New("chain ID cannot be an empty string")
	}
	logger.Info("registering HTLC manager", "chain_id", cid)
	_, exists := s.Mgrs[cid]
	if exists {
		return errors.New("manager already registered")
	}

	s.Mgrs[cid] = mgr
	return nil
}

func (s *HTLCSwitch) Manager(chainID string) (HTLCManager) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.Mgrs[chainID]
}