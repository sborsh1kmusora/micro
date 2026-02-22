package item

import (
	"context"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

func (s *service) Create(ctx context.Context, itemInfo *model.ItemInfo) (string, error) {
	uuid, err := s.itemRepo.Create(ctx, itemInfo)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
