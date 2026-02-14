package grpc

import (
	"context"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

type InventoryClient interface {
	ListItems(context.Context, []string) ([]*model.Item, error)
}

type PaymentClient interface {
	PayOrder(context.Context, string, string, string) (string, error)
}
