package item

import (
	"context"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/converter"
)

func (r *repo) List(_ context.Context) ([]*model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Item, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, converter.ItemToModel(item))
	}

	return result, nil
}
