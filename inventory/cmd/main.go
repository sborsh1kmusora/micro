package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

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
	if req == nil || req.Filter == nil {
		return &inventoryV1.ListItemsResponse{
			Items: mapToSlice(s.items),
		}, nil
	}

	filter := req.Filter

	if len(filter.Uuids) == 0 {
		return &inventoryV1.ListItemsResponse{
			Items: mapToSlice(s.items),
		}, nil
	}

	var result []*inventoryV1.Item

	if len(filter.Uuids) > 0 {
		uuidSet := make(map[string]struct{}, len(filter.Uuids))
		for _, n := range filter.Uuids {
			uuidSet[n] = struct{}{}
		}

		for _, item := range s.items {
			if _, ok := uuidSet[item.GetUuid()]; ok {
				result = append(result, item)
			}
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

	s := grpc.NewServer()

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
