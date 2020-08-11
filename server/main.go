package main

import (
	"context"
	"flag"
	"fmt"
	podlifecycle "github.com/jwenz723/podlifecycle/server/proto"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	grpcAddr := flag.String("grpcAddr", ":8080", "address to expose grpc on")
	requestDuration := flag.Duration("requestDuration", 0, "duration that each grpc request should take to process")
	flag.Parse()
	fmt.Println(*grpcAddr, *requestDuration)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var g run.Group
	{
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Error("failed to start grpc listener", zap.Error(err))
		}
		service := stufferService{
			l:               logger,
			requestDuration: *requestDuration,
		}
		grpcServer := NewGRPCServerFromListener(lis)

		g.Add(func() error {
			podlifecycle.RegisterStufferServer(grpcServer.Server(), &service)
			logger.Info("starting grpc server...", zap.String("addr", *grpcAddr))
			return grpcServer.Start()
		}, func(err error) {
			logger.Info("shutting down grpc server...")
			grpcServer.Stop()
			lis.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				logger.Info("received signal", zap.String("signal", sig.String()))
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				logger.Info("cancel interrupt")
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Info("starting...")
	logger.Info("exiting", zap.Error(g.Run()))
}

var _ podlifecycle.StufferServer = &stufferService{}

type stufferService struct {
	l               *zap.Logger
	requestDuration time.Duration
}

func (s *stufferService) DoStuff(_ context.Context, req *podlifecycle.StuffRequest) (*podlifecycle.StuffResponse, error) {
	s.l.Info("DoStuff invoked3", zap.String("name", req.Name))
	time.Sleep(s.requestDuration)
	s.l.Info("DoStuff completed3", zap.String("name", req.Name))
	return &podlifecycle.StuffResponse{Name: req.Name}, nil
}
