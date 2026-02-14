package order

import (
	gen "github.com/sborsh1kmusora/micro/order/internal/client/grpc"
	"github.com/sborsh1kmusora/micro/order/internal/repository"
	def "github.com/sborsh1kmusora/micro/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo repository.OrderRepository

	inventoryCl gen.InventoryClient
	paymentCL   gen.PaymentClient
}

func NewService(
	orderRepo repository.OrderRepository,
	inventoryCl gen.InventoryClient,
	paymentCL gen.PaymentClient,
) *service {
	return &service{
		orderRepo:   orderRepo,
		inventoryCl: inventoryCl,
		paymentCL:   paymentCL,
	}
}
