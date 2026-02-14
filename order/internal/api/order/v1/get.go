package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/sborsh1kmusora/micro/order/internal/converter"
	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(
	ctx context.Context,
	params orderV1.GetOrderParams,
) (orderV1.GetOrderRes, error) {
	order, err := a.orderSvc.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Message: "Order not found",
				Code:    http.StatusNotFound,
			}, nil
		}
		return &orderV1.InternalServerError{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	res := converter.OrderModelToOpenAPI(order)

	return res, nil
}
