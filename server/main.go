package main

import (
	"context"
	"flag"
	"fmt"
	example "github.com/jwenz723/podlifecycle/server/proto"
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
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var g run.Group
	{
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Error("failed to start grpc listener", zap.Error(err))
		}
		service := exampleService{
			l: logger,
		}
		grpcServer := NewGRPCServerFromListener(lis)

		g.Add(func() error {
			example.RegisterExampleServer(grpcServer.Server(), &service)
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

	logger.Info("exiting", zap.Error(g.Run()))
}

type exampleService struct {
	l               *zap.Logger
	requestDuration time.Duration
}

// Work implements example.ExampleServer
func (s *exampleService) Work(_ context.Context, req *example.WorkItem) (*example.WorkResponse, error) {
	s.l.Info("Work invoked", zap.String("name", req.Name), zap.Int("size", int(req.Size)))

	// sleep based upon the specified request size to simulate slow requests
	time.Sleep(time.Duration(req.Size) * time.Second)

	return &example.WorkResponse{Name: req.Name}, nil
}
