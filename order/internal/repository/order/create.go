package order

import (
	"context"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	"github.com/sborsh1kmusora/micro/order/internal/repository/converter"
)

func (r *repo) Create(_ context.Context, uuid string, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[uuid] = converter.OrderModelToRepo(order)

	return nil
}
