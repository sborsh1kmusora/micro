package item

import (
	"context"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

func (s *service) Create(ctx context.Context, item *model.ItemInfo) (string, error) {
	uuid, err := s.itemRepo.Create(ctx, item)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
