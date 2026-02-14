package item

import (
	"context"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro/inventory/internal/repository/model"
)

func (r *repo) Create(_ context.Context, info *model.ItemInfo) (string, error) {
	newUUID := uuid.NewString()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[newUUID] = &repoModel.Item{
		UUID: newUUID,
		Info: converter.ItemInfoToRepoModel(info),
	}

	return newUUID, nil
}
