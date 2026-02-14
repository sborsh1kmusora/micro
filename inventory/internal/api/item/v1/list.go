package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/inventory/internal/converter"
	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func (a *api) ListItems(
	ctx context.Context,
	req *inventoryV1.ListItemsRequest,
) (*inventoryV1.ListItemsResponse, error) {
	items, err := a.itemService.List(ctx, req.Uuids)
	if err != nil {
		if errors.Is(err, model.ErrItemNotFound) {
			return nil, status.Error(codes.NotFound, "item not found")
		}
		return nil, err
	}

	return &inventoryV1.ListItemsResponse{
		Items: converter.ListItemsToProto(items),
	}, nil
}
