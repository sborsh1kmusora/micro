package repository

import (
	"context"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

type ItemRepository interface {
	Create(context.Context, *model.ItemInfo) (string, error)
	Get(context.Context, string) (*model.Item, error)
	List(context.Context) ([]*model.Item, error)
}
