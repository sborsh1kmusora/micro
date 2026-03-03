package item

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sborsh1kmusora/micro/inventory/internal/repository"
)

var _ repository.ItemRepository = (*repo)(nil)

type repo struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) *repo {
	collection := client.Database("inventory").Collection("items")

	return &repo{
		collection: collection,
	}
}
