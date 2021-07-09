package util

import (
	log "github.com/inconshreveable/log15"
	"fmt"
)

var logger = log.New("module", "grpc")

type grpcLog struct {

}

func (*grpcLog) Info(args ...interface{}) {
	logger.Info("GRPC message", args)
}

func (*grpcLog) Infoln(args ...interface{}) {
	logger.Info("GRPC message\n", args)
}

func (*grpcLog) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args)
	logger.Info(msg)
}

func (*grpcLog) Warning(args ...interface{}) {
	logger.Warn("GRPC message", args)
}

func (*grpcLog) Warningln(args ...interface{}) {
	logger.Warn("GRPC message\n", args)
}

func (*grpcLog) Warningf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args)
	logger.Warn(msg)
}

func (*grpcLog) Error(args ...interface{}) {
	logger.Error("GRPC message", args)
}

func (*grpcLog) Errorln(args ...interface{}) {
	logger.Error("GRPC message\n", args)
}

func (*grpcLog) Errorf(format string, args ...interface{}) {
	logger.Error(format, args)
}

func (*grpcLog) Fatal(args ...interface{}) {
	logger.Error("GRPC message", args)
	panic(args)
}

func (*grpcLog) Fatalln(args ...interface{}) {
	logger.Error("GRPC message\n", args)
	panic(args)
}

func (*grpcLog) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args)
	logger.Error(msg)
	panic(msg)
}

func (*grpcLog) V(l int) bool {
	return true
}



