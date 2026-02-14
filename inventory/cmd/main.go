package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	itemV1API "github.com/sborsh1kmusora/micro/inventory/internal/api/item/v1"
	"github.com/sborsh1kmusora/micro/inventory/internal/interceptor"
	itemRepo "github.com/sborsh1kmusora/micro/inventory/internal/repository/item"
	itemService "github.com/sborsh1kmusora/micro/inventory/internal/service/item"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	repo := itemRepo.NewRepository()
	svc := itemService.NewService(repo)
	api := itemV1API.NewAPI(svc)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on port %d", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		log.Printf("HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down gRPC server...")

	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("HTTP server stopped")
	}

	s.GracefulStop()

	log.Println("Server gracefully stopped")
}
