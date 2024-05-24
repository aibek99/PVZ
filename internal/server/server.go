package server

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"Homework-1/internal/cache"
	"Homework-1/internal/config"
	"Homework-1/internal/database"
	"Homework-1/internal/kafka"
)

// Server is
type Server struct {
	config     *config.Config
	gRPC       *grpc.Server
	dataStore  database.Datastore
	producer   kafka.Producer
	cacheStore cache.Store
}

// NewServer is
func NewServer(
	cfg *config.Config,
	dataStore database.Datastore,
	cacheStore cache.Store,
	producer kafka.Producer,
) *Server {
	server := &Server{
		config:     cfg,
		dataStore:  dataStore,
		cacheStore: cacheStore,
		producer:   producer,
	}
	return server
}

// Run is
func (s *Server) Run(ctx context.Context) error {
	lisSecure, err := net.Listen("tcp", s.config.Server.GRPCPort)
	if err != nil {
		return err
	}
	defer lisSecure.Close()
	creds, err := credentials.NewServerTLSFromFile(s.config.TLS.CertPath, s.config.TLS.KeyPath)
	if err != nil {
		return err
	}

	s.gRPC = grpc.NewServer(grpc.Creds(creds), CombinedInterceptor(s.producer))
	reflection.Register(s.gRPC)
	s.mapHandlers()

	var wg sync.WaitGroup

	// Starting the GRPC server in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[server][Run] GRPC server is started on %s\n", s.config.Server.GRPCPort)
		if err := s.gRPC.Serve(lisSecure); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[server][Run] HTTP server is started on %s\n", s.config.Server.HTTPPort)
		if err := s.StartGatewayRouter(ctx); err != nil {
			log.Printf("failed to start gateway router: %v", err)
		}
	}()
	<-ctx.Done()

	// Shutdown the HTTPS server
	s.gRPC.GracefulStop()

	log.Println("[server][Run] Secure server shutdown gracefully")

	wg.Wait()

	return nil
}
