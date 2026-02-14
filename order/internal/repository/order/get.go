package order

import (
	"context"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	"github.com/sborsh1kmusora/micro/order/internal/repository/converter"
)

func (r *repo) Get(_ context.Context, uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}
