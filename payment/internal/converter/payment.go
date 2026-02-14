package converter

import (
	"github.com/sborsh1kmusora/micro/payment/internal/model"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(req *paymentV1.PayOrderRequest) *model.PayOrderRequest {
	return &model.PayOrderRequest{
		OrderUUID:     req.OrderUuid,
		UserUUID:      req.UserUuid,
		PaymentMethod: req.PaymentMethod.String(),
	}
}
