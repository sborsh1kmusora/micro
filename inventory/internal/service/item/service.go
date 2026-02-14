package item

import (
	"github.com/sborsh1kmusora/micro/inventory/internal/repository"
	def "github.com/sborsh1kmusora/micro/inventory/internal/service"
)

var _ def.ItemService = (*service)(nil)

type service struct {
	itemRepo repository.ItemRepository
}

func NewService(itemRepo repository.ItemRepository) *service {
	return &service{itemRepo: itemRepo}
}
