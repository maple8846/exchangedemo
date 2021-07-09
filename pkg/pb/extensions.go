package pb

import (
	"golang.org/x/crypto/sha3"
	"hash"
	"encoding/binary"
)



func (o *OrderV1) Hash() (res [32]byte) {
	quantity := make([]byte, 8)
	binary.LittleEndian.PutUint64(quantity, uint64(o.Quantity))


	price := make([]byte, 8)
        binary.LittleEndian.PutUint64(price, uint64(o.Price))
	h := sha3.New256()
	guardedWrite(h, []byte(o.OrderId))
	guardedWrite(h, []byte(o.UserPubKey))
	guardedWrite(h, []byte(o.MakerAsset))
	guardedWrite(h, []byte(o.TakerAsset))
	guardedWrite(h, []byte(quantity))
	guardedWrite(h, []byte(price))
	guardedWrite(h, []byte(o.CreatedAt))
	guardedWrite(h, []byte(o.Memo))
	guardedWrite(h, []byte(o.CreatedAt))
	copy(res[:], h.Sum(nil))
	return res
}

func guardedWrite(h hash.Hash, b []byte) {
	_, err := h.Write(b)
	if err != nil {
		panic(err)
	}
}
