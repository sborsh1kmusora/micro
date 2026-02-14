package order

import (
	"context"
	"errors"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.ErrOrderCantBeCancelled
	}

	order.Status = model.OrderStatusCancelled

	if err := s.orderRepo.Create(ctx, uuid, order); err != nil {
		return err
	}

	return nil
}
