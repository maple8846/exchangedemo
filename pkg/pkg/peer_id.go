package pkg

import (
	"github.com/btcsuite/btcd/btcec"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
	"github.com/pkg/errors"
	"strconv"
)

type PeerID struct {
	Pubkey *btcec.PublicKey
	Addr   string
	Port   int
}

func (p *PeerID) String() string {
	return fmt.Sprintf("%s::%s::%d", hexutil.Encode(p.Pubkey.SerializeCompressed()), p.Addr, p.Port)
}

func PeerIDFromStr(str string) (*PeerID, error) {
	parts := strings.Split(str, "::")

	if len(parts) != 3 {
		return nil, errors.New("invalid peer ID")
	}

	var pub *btcec.PublicKey
	var addr string
	var port int

	pubBytes, err := hexutil.Decode(parts[0])
	if err != nil {
		return nil, err
	}
	pub, err = btcec.ParsePubKey(pubBytes, btcec.S256())
	if err != nil {
		return nil, err
	}

	if parts[1] == "" {
		return nil, errors.New("invalid address in peer ID")
	}
	addr = parts[1]

	port, err = strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.New("invalid port")
	}

	return &PeerID{pub, addr, port}, nil
}
