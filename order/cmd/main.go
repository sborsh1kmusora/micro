package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/sborsh1kmusora/micro/order/internal/api/order/v1"
	inventoryV1Client "github.com/sborsh1kmusora/micro/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/sborsh1kmusora/micro/order/internal/client/grpc/payment/v1"
	customMiddleware "github.com/sborsh1kmusora/micro/order/internal/middleware"
	orderRepo "github.com/sborsh1kmusora/micro/order/internal/repository/order"
	orderSvc "github.com/sborsh1kmusora/micro/order/internal/service/order"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	inventoryGrpcClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentGrpcClient := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryCl := inventoryV1Client.NewClient(inventoryGrpcClient)
	paymentCl := paymentV1Client.NewClient(paymentGrpcClient)

	repo := orderRepo.NewRepository()
	svc := orderSvc.NewService(repo, inventoryCl, paymentCl)

	api := v1.NewAPI(svc)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Error creating order server: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:        net.JoinHostPort("localhost", httpPort),
		Handler:     r,
		ReadTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Starting server on port %s", httpPort)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
}
