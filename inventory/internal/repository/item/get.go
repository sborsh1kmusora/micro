package item

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro/inventory/internal/repository/model"
)

func (r *repo) Get(ctx context.Context, uuid string) (*model.Item, error) {
	var item repoModel.Item

	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrItemNotFound
		}

		return nil, err
	}

	return converter.ItemToModel(&item), nil
}
