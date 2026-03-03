package item

import (
	"context"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro/inventory/internal/repository/model"
)

func (r *repo) Create(ctx context.Context, info *model.ItemInfo) (string, error) {
	newUUID := uuid.NewString()

	item := &repoModel.Item{
		UUID: newUUID,
		Info: converter.ItemInfoToRepoModel(info),
	}

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return "", err
	}

	return newUUID, nil
}
