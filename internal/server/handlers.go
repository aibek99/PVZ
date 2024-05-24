package server

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	BoxDelivery "Homework-1/internal/box/delivery"
	BoxUseCase "Homework-1/internal/box/usecase"
	"Homework-1/internal/kafka"
	"Homework-1/internal/metrics"
	kafkaModel "Homework-1/internal/model/kafka"
	OrderDelivery "Homework-1/internal/order/delivery"
	OrderUseCase "Homework-1/internal/order/usecase"
	PVZDelivery "Homework-1/internal/pvz/delivery"
	PVZUseCase "Homework-1/internal/pvz/usecase"
	"Homework-1/pkg/api/box_v1"
	"Homework-1/pkg/api/order_v1"
	"Homework-1/pkg/api/pvz_v1"
)

// CombinedInterceptor is
func CombinedInterceptor(kafkaProducer kafka.Producer) grpc.ServerOption {
	return grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		authInterceptor(),               // Authentication interceptor
		KafkaInterceptor(kafkaProducer), // kafka interceptor
		MetricsInterceptor(),            // metrics interceptor
	))
}

func authInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// BasicAuthInterceptor logic

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) < 1 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		const prefix = "Basic "
		if !strings.HasPrefix(authHeader[0], prefix) {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not Basic")
		}

		encodedCredentials := authHeader[0][len(prefix):]
		credentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to decode authorization token")
		}

		parts := strings.SplitN(string(credentials), ":", 2)
		if len(parts) != 2 {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		username := parts[0]
		password := parts[1]

		// Check the username and password against your expected values
		if username != "Homework_3" || password != "test" {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return handler(ctx, request)
	}
}

// MetricsInterceptor is
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()
		resp, err := handler(ctx, request)
		metrics.ResponseTimes.Observe(time.Since(startTime).Seconds())
		metrics.TotalRequests.Inc()

		st, ok := status.FromError(err)
		if !ok {
			log.Printf("unable to get error from status")
		}
		metrics.ResponseStatus.WithLabelValues(st.Code().String()).Inc()

		if err != nil {
			metrics.ErrorRates.Inc()
		}

		return resp, err
	}

}

// KafkaInterceptor is
func KafkaInterceptor(producer kafka.Producer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// KafkaInterceptor logic
		requestData, err := json.Marshal(request)
		if err != nil {
			log.Printf("Failed to marshal request to JSON: %v", err)
			return nil, err
		}
		// Created Kafka Message
		err = producer.SendMessage(producer.Topic(), kafkaModel.Message{
			Method:    info.FullMethod,
			Request:   string(requestData),
			Timestamp: time.Now().Round(time.Minute),
		})
		if err != nil {
			log.Printf("Failed to send Kafka message: %v", err)
		}

		return handler(ctx, request)
	}
}

func (s *Server) mapHandlers() {
	orderUseCase := OrderUseCase.NewOrderUseCase(s.dataStore, s.cacheStore)
	orderHandlers := OrderDelivery.NewOrdersHandler(orderUseCase)
	order_v1.RegisterOrderServiceServer(s.gRPC, orderHandlers)

	boxUseCase := BoxUseCase.NewBoxUseCase(s.dataStore, s.cacheStore)
	boxHandlers := BoxDelivery.NewBoxHandler(boxUseCase)
	box_v1.RegisterBoxServiceServer(s.gRPC, boxHandlers)

	pvzUseCase := PVZUseCase.NewPVZUseCase(s.dataStore, s.cacheStore)
	pvzHandlers := PVZDelivery.NewPVZHandler(pvzUseCase)
	pvz_v1.RegisterPVZServiceServer(s.gRPC, pvzHandlers)
}

// StartGatewayRouter is
func (s *Server) StartGatewayRouter(ctx context.Context) error {
	certPool := x509.NewCertPool()

	serverCert, err := os.ReadFile(s.config.TLS.CACrt)
	if err != nil {
		return err
	}

	if !certPool.AppendCertsFromPEM(serverCert) {
		return err
	}

	log.Println("Setting up client credentials...")
	creds := credentials.NewClientTLSFromCert(certPool, "")
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	mux := runtime.NewServeMux()

	log.Printf("Registering Box service handler with endpoint %s...", s.config.Server.GRPCPort)
	err = box_v1.RegisterBoxServiceHandlerFromEndpoint(ctx, mux, s.config.Server.GRPCPort, opts)
	if err != nil {
		return err
	}

	err = pvz_v1.RegisterPVZServiceHandlerFromEndpoint(ctx, mux, s.config.Server.GRPCPort, opts)
	if err != nil {
		return err
	}

	err = order_v1.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, s.config.Server.GRPCPort, opts)
	if err != nil {
		return err
	}
	// Start HTTP server (and proxy calls to gRPC server)
	log.Printf("Starting HTTP server on port %v", s.config.Server.HTTPPort)

	server := &http.Server{
		Addr:         s.config.Server.HTTPPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return server.ListenAndServe()
}
