package item

import (
	"context"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

func (s *service) List(ctx context.Context, uuids []string) ([]*model.Item, error) {
	if len(uuids) > 0 {
		result := make([]*model.Item, 0, len(uuids))
		for _, uuid := range uuids {
			item, err := s.itemRepo.Get(ctx, uuid)
			if err != nil {
				return nil, err
			}
			result = append(result, item)
		}

		return result, nil
	}

	return s.itemRepo.List(ctx)
}
