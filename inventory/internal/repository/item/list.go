package item

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro/inventory/internal/repository/model"
)

func (r *repo) List(ctx context.Context) ([]*model.Item, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var notes []*repoModel.Item
	err = cursor.All(ctx, &notes)
	if err != nil {
		return nil, err
	}

	var result []*model.Item
	for _, item := range notes {
		result = append(result, converter.ItemToModel(item))
	}

	return result, nil
}
