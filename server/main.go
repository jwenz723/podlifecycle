package main

import (
	"context"
	"flag"
	"fmt"
	podlifecycle "github.com/jwenz723/podlifecycle/proto"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	requestDuration := flag.Duration("requestDuration", 0, "duration that each grpc request should take to process")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var g run.Group
	{
		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			logger.Error("failed to start grpc listener", zap.Error(err))
		}
		service := stufferService{
			l: logger,
			requestDuration: *requestDuration,
		}
		server := grpc.NewServer()
		g.Add(func() error {
			podlifecycle.RegisterStufferServer(server, &service)
			return server.Serve(lis)
		}, func(err error) {
			logger.Info("shutting down grpc server")
			server.GracefulStop()
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
	l *zap.Logger
	requestDuration time.Duration
}

func (s *stufferService) DoStuff(ctx context.Context, req *podlifecycle.StuffRequest) (*podlifecycle.StuffResponse, error) {
	s.l.Info("DoStuff invoked2", zap.String("name", req.Name))
	time.Sleep(s.requestDuration)
	s.l.Info("DoStuff completed2", zap.String("name", req.Name))
	return &podlifecycle.StuffResponse{Name: req.Name}, nil
}

