package repository

import (
	"context"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

type OrderRepository interface {
	Create(context.Context, string, *model.Order) error
	Get(context.Context, string) (*model.Order, error)
}
