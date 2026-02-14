package order

import (
	"context"
	"errors"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Order, error) {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}
