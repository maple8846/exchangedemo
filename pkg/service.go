package pkg

import "context"

type Service interface {
	Start() error
	Stop() error
}

type BaseService struct {
	Ctx context.Context
}