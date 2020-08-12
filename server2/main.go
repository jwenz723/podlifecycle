package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version = 1

func main() {
	httpAddr := flag.String("httpAddr", ":8080", "address to bind http server to")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	logger = logger.With(zap.Int("version", version))

	var g run.Group
	{
		httpServer := NewHTTPService(httpServiceConfig{}, logger)
		l, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			panic(err)
		}
		g.Add(func() error {
			logger.Info("starting HTTP server", zap.String("addr", *httpAddr))
			return httpServer.Start(l)
		}, func(err error) {
			logger.Info("starting shutdown of http")
			httpServer.Stop(context.Background())
			l.Close()
			logger.Info("completed shutdown of http")
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

type httpServiceConfig struct {
	PreStopWaitTime int `json:"preStopWaitTime"`
}

type httpService struct {
	config httpServiceConfig
	logger *zap.Logger
	server http.Server
}

func NewHTTPService(config httpServiceConfig, logger *zap.Logger) httpService {
	return httpService{
		config: config,
		logger: logger,
		server: http.Server{},
	}
}

func (h *httpService) Start(l net.Listener) error {
	m := mux.NewRouter()
	m.Path("/blocker").HandlerFunc(h.handleBlocker())
	m.Path("/configure").HandlerFunc(h.handleConfigure())
	m.Path("/liveness").HandlerFunc(h.handleLiveness())
	m.Path("/prestop").HandlerFunc(h.handlePreStop())
	m.Path("/readiness").HandlerFunc(h.handleReadiness())
	h.server.Handler = m
	return h.server.Serve(l)
}

func (h *httpService) Stop(ctx context.Context) {
	h.server.Shutdown(ctx)
}

func (h *httpService) handleLiveness() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h.logger.Info("liveness invoked")
		writer.WriteHeader(http.StatusOK)
	}
}

func (h *httpService) handlePreStop() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h.logger.Info("prestop invoked")
		time.Sleep(time.Duration(h.config.PreStopWaitTime) * time.Second)
		writer.WriteHeader(http.StatusOK)
		h.logger.Info("prestop completed")
	}
}

func (h *httpService) handleReadiness() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h.logger.Info("readiness invoked")
		writer.WriteHeader(http.StatusOK)
	}
}

func (h *httpService) handleConfigure() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var c httpServiceConfig
		err := json.NewDecoder(request.Body).Decode(&c)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Info("configure invoked", zap.Int("preStopWaitTime", c.PreStopWaitTime))
		h.config = c
		writer.WriteHeader(http.StatusOK)
	}
}

func (h *httpService) handleBlocker() http.HandlerFunc {
	type blockerReq struct {
		Seconds int `json:"seconds""`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		var c blockerReq
		err := json.NewDecoder(request.Body).Decode(&c)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Info("blocker invoked", zap.Int("seconds", c.Seconds))
		time.Sleep(time.Duration(c.Seconds) * time.Second)
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "ok")
		h.logger.Info("blocker completed", zap.Int("seconds", c.Seconds))
	}
}
