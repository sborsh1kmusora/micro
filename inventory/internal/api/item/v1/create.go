package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/inventory/internal/converter"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func (a *api) CreateItem(
	ctx context.Context,
	req *inventoryV1.CreateItemRequest,
) (*inventoryV1.CreateItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	itemUUID, err := a.itemService.Create(ctx, converter.InfoProtoToModel(req.Info))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal Server Error")
	}

	return &inventoryV1.CreateItemResponse{
		Uuid: itemUUID,
	}, nil
}
