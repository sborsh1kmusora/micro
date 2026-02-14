package converter

import (
	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func OrderModelToOpenAPI(order *model.Order) *orderV1.Order {
	return &orderV1.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		ItemUuids:       order.ItemUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: orderV1.NewOptNilString(order.TransactionUUID),
		PaymentMethod:   orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(order.PaymentMethod)),
		Status:          orderV1.OrderStatus(order.Status),
	}
}
