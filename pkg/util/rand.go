package util

import "crypto/rand"

func Rand32() ([32]byte) {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return b
}