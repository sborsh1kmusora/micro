package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/sborsh1kmusora/micro/order/internal/middleware"
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

	pendingPaymentStatus = "PENDING_PAYMENT"
	paidStatus           = "PAID"
	cancelledStatus      = "CANCELLED"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) SaveOrder(_ context.Context, uuid string, order *orderV1.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = order

	return nil
}

func (s *OrderStorage) GetOrder(_ context.Context, uuid string) (*orderV1.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

type OrderHandler struct {
	OrderStorage    *OrderStorage
	InventoryClient inventoryV1.InventoryServiceClient
	PaymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(
	storage *OrderStorage,
	inventoryCl inventoryV1.InventoryServiceClient,
	paymentCl paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		OrderStorage:    storage,
		InventoryClient: inventoryCl,
		PaymentClient:   paymentCl,
	}
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	resp, err := h.InventoryClient.ListItems(ctx, &inventoryV1.ListItemsRequest{
		Uuids: req.ItemUuids,
	})
	if err != nil {
		log.Printf("Error listing items: %v", err)
		return &orderV1.BadRequestError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}, nil
	}

	totalPrice := 0.0
	for _, item := range resp.Items {
		totalPrice += item.Info.Price
	}

	orderUUID := uuid.NewString()

	order := &orderV1.Order{
		UUID:       orderUUID,
		UserUUID:   req.UserUUID,
		ItemUuids:  req.ItemUuids,
		TotalPrice: totalPrice,
		Status:     pendingPaymentStatus,
	}

	if err := h.OrderStorage.SaveOrder(ctx, orderUUID, order); err != nil {
		log.Printf("Error creating order %v with uuid %s: %v", order, orderUUID, err)
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(
	ctx context.Context,
	params orderV1.GetOrderParams,
) (orderV1.GetOrderRes, error) {
	order, err := h.OrderStorage.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		if errors.Is(err, ErrOrderNotFound) {
			log.Printf("Order with uuid %s not found", params.OrderUUID)
			return &orderV1.NotFoundError{
				Message: "Order not found",
				Code:    http.StatusNotFound,
			}, nil
		}

		return &orderV1.InternalServerError{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	order, err := h.OrderStorage.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		if errors.Is(err, ErrOrderNotFound) {
			log.Printf("Order with uuid %s not found", params.OrderUUID)
			return &orderV1.NotFoundError{
				Message: "Order not found",
				Code:    http.StatusNotFound,
			}, nil
		}

		return &orderV1.InternalServerError{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	if order.Status == paidStatus {
		return &orderV1.ConflictError{
			Message: "Order has already paid",
			Code:    http.StatusConflict,
		}, nil
	}

	resp, err := h.PaymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentProtoFromAPI(req.PaymentMethod),
	})
	if err != nil {
		log.Printf("Error paying order %s by user %s with %s: %v", params.OrderUUID, order.UserUUID, req.PaymentMethod, err)
		return &orderV1.InternalServerError{
			Message: "Failed to pay order",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	order.Status = paidStatus
	order.PaymentMethod = orderV1.NewOptPaymentMethod(req.PaymentMethod)
	order.TransactionUUID = orderV1.NewOptNilString(resp.TransactionUuid)

	if err := h.OrderStorage.SaveOrder(ctx, order.UUID, order); err != nil {
		log.Printf("Error saving order %v with uuid %s: %v", order, order.UUID, err)
		return &orderV1.InternalServerError{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: resp.TransactionUuid,
	}, nil
}

func paymentProtoFromAPI(
	m orderV1.PaymentMethod,
) paymentV1.PaymentMethod {
	switch m {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func (h *OrderHandler) CancelOrder(
	ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	order, err := h.OrderStorage.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		if errors.Is(err, ErrOrderNotFound) {
			log.Printf("Order with uuid %s not found", params.OrderUUID)
			return &orderV1.NotFoundError{
				Message: "Order not found",
				Code:    http.StatusNotFound,
			}, nil
		}
	}

	if order.Status == paidStatus || order.Status == cancelledStatus {
		return &orderV1.ConflictError{
			Message: "Order has already paid or cancelled",
			Code:    http.StatusConflict,
		}, nil
	}

	order.Status = cancelledStatus

	if err := h.OrderStorage.SaveOrder(ctx, order.UUID, order); err != nil {
		log.Printf("Error saving order %v with uuid %s: %v", order, order.UUID, err)
		return &orderV1.InternalServerError{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return &orderV1.CancelOrderNoContent{}, nil
}

func main() {
	storage := NewOrderStorage()

	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryCl := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentCl := paymentV1.NewPaymentServiceClient(paymentConn)

	handler := NewOrderHandler(storage, inventoryCl, paymentCl)

	orderServer, err := orderV1.NewServer(handler)
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
