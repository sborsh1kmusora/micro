package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	transactionUUID, err := a.orderSvc.Pay(ctx, params.OrderUUID, string(req.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Message: "Order not found",
				Code:    http.StatusNotFound,
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaidOrCancelled):
			return &orderV1.ConflictError{
				Message: "Order has already paid or cancelled",
				Code:    http.StatusConflict,
			}, nil
		default:
			return &orderV1.InternalServerError{
				Message: "Internal Server Error",
				Code:    http.StatusInternalServerError,
			}, nil
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
