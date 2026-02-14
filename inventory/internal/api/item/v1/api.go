package v1

import (
	"github.com/sborsh1kmusora/micro/inventory/internal/service"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	itemService service.ItemService
}

func NewAPI(itemService service.ItemService) *api {
	return &api{
		itemService: itemService,
	}
}
