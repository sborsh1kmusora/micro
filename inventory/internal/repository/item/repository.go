package item

import (
	"sync"

	"github.com/sborsh1kmusora/micro/inventory/internal/repository"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/model"
)

var _ repository.ItemRepository = (*repo)(nil)

type repo struct {
	mu    sync.RWMutex
	items map[string]*model.Item
}

func NewRepository() *repo {
	return &repo{
		items: make(map[string]*model.Item),
	}
}
