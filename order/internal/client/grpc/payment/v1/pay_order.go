package v1

import (
	"context"

	"github.com/sborsh1kmusora/micro/order/internal/client/converter"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(
	ctx context.Context,
	orderUUID, userUUID,
	paymentMethod string,
) (string, error) {
	resp, err := c.genClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: converter.ToProtoPaymentMethod(paymentMethod),
	})
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}
