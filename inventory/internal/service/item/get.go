package item

import (
	"context"
	"errors"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Item, error) {
	item, err := s.itemRepo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrItemNotFound) {
			return nil, model.ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}
