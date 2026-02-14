package service

import (
	"context"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

type OrderService interface {
	Create(context.Context, string, []string) (string, float64, error)
	Get(context.Context, string) (*model.Order, error)
	Pay(context.Context, string, string) (string, error)
	Cancel(context.Context, string) error
}
