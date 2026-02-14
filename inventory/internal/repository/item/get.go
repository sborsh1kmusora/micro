package item

import (
	"context"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/converter"
)

func (r *repo) Get(_ context.Context, uuid string) (*model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[uuid]
	if !ok {
		return nil, model.ErrItemNotFound
	}

	return converter.ItemToModel(item), nil
}
