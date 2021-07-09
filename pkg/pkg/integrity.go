package pkg

type Hasher interface {
	Hash() [32]byte
}