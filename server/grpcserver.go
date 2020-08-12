package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type GRPCServer struct {
	// Listen address for the server specified as hostname:port
	address string
	// Listener for handling network requests
	listener net.Listener
	// GRPC server
	server *grpc.Server
	// Server for gRPC Health Check Protocol.
	healthServer *health.Server
}

// NewGRPCServerFromListener creates a new implementation of a GRPCServer given
// an existing net.Listener instance using default keepalive
func NewGRPCServerFromListener(listener net.Listener) *GRPCServer {
	grpcServer := &GRPCServer{
		address:  listener.Addr().String(),
		listener: listener,
	}

	grpcServer.server = grpc.NewServer()
	grpcServer.healthServer = health.NewServer()
	healthpb.RegisterHealthServer(grpcServer.server, grpcServer.healthServer)

	return grpcServer
}

// Start starts the underlying grpc.Server
func (gServer *GRPCServer) Start() error {
	// if health check is enabled, set the health status for all registered services
	if gServer.healthServer != nil {
		for name := range gServer.server.GetServiceInfo() {
			gServer.healthServer.SetServingStatus(
				name,
				healthpb.HealthCheckResponse_SERVING,
			)
		}

		gServer.healthServer.SetServingStatus(
			"",
			healthpb.HealthCheckResponse_SERVING,
		)
	}
	return gServer.server.Serve(gServer.listener)
}

func (gServer *GRPCServer) Stop() {
	gServer.server.GracefulStop()
}

// Server returns the grpc.Server for the GRPCServer instance
func (gServer *GRPCServer) Server() *grpc.Server {
	return gServer.server
}
