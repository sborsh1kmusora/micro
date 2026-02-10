package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/inventory/internal/interceptor"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

var ErrItemNotFound = errors.New("item not found")

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	items map[string]*inventoryV1.Item
}

func (s *inventoryService) GetItem(
	_ context.Context,
	req *inventoryV1.GetItemRequest,
) (*inventoryV1.GetItemResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "item with uuid %s not found", req.GetUuid())
	}

	return &inventoryV1.GetItemResponse{
		Item: item,
	}, nil
}

func (s *inventoryService) ListItems(
	_ context.Context,
	req *inventoryV1.ListItemsRequest,
) (*inventoryV1.ListItemsResponse, error) {
	if req == nil || req.GetUuids() == nil || len(req.GetUuids()) == 0 {
		return &inventoryV1.ListItemsResponse{
			Items: mapToSlice(s.items),
		}, nil
	}

	uuids := req.GetUuids()

	var result []*inventoryV1.Item

	uuidSet := make(map[string]*inventoryV1.Item, len(s.items))
	for u, i := range s.items {
		uuidSet[u] = i
	}

	for _, u := range uuids {
		if _, ok := uuidSet[u]; ok {
			result = append(result, uuidSet[u])
		} else {
			log.Printf("item with uuid %s not found", u)
			return nil, ErrItemNotFound
		}
	}

	return &inventoryV1.ListItemsResponse{
		Items: result,
	}, nil
}

func (s *inventoryService) CreateItem(
	_ context.Context,
	req *inventoryV1.CreateItemRequest,
) (*inventoryV1.CreateItemResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	newUuid := uuid.NewString()
	item := &inventoryV1.Item{
		Uuid: newUuid,
		Info: req.GetInfo(),
	}

	s.items[newUuid] = item

	return &inventoryV1.CreateItemResponse{
		Uuid: newUuid,
	}, nil
}

func mapToSlice(m map[string]*inventoryV1.Item) []*inventoryV1.Item {
	result := make([]*inventoryV1.Item, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	service := &inventoryService{
		items: make(map[string]*inventoryV1.Item),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
