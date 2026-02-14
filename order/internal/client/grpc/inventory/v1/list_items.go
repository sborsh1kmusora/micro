package v1

import (
	"context"

	clientConverter "github.com/sborsh1kmusora/micro/order/internal/client/converter"
	"github.com/sborsh1kmusora/micro/order/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func (c *client) ListItems(ctx context.Context, uuids []string) ([]*model.Item, error) {
	items, err := c.genClient.ListItems(ctx, &inventoryV1.ListItemsRequest{
		Uuids: uuids,
	})
	if err != nil {
		return nil, err
	}

	return clientConverter.ItemListToModel(items), nil
}
