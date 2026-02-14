package v1

import (
	def "github.com/sborsh1kmusora/micro/order/internal/client/grpc"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	genClient inventoryV1.InventoryServiceClient
}

func NewClient(genClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		genClient: genClient,
	}
}
