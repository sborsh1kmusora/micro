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

func (a *api) GetItem(ctx context.Context, getItemReq *inventoryV1.GetItemRequest) (*inventoryV1.GetItemResponse, error) {
	item, err := a.itemService.Get(ctx, getItemReq.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrItemNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return &inventoryV1.GetItemResponse{
		Item: converter.ItemModelToProto(item),
	}, nil
}
