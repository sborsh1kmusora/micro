package order

import (
	"sync"

	"github.com/sborsh1kmusora/micro/order/internal/repository"
	"github.com/sborsh1kmusora/micro/order/internal/repository/model"
)

var _ repository.OrderRepository = (*repo)(nil)

type repo struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewRepository() *repo {
	return &repo{
		orders: make(map[string]*model.Order),
	}
}
