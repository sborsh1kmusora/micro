package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	orderUUID, totalPrice, err := a.orderSvc.Create(ctx, req.UserUUID, req.ItemUuids)
	if err != nil {
		if errors.Is(err, model.ErrItemNotFound) {
			return &orderV1.BadRequestError{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
