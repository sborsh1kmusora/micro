package converter

import (
	"time"

	"github.com/samber/lo"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

func ItemListToModel(resp *inventoryV1.ListItemsResponse) []*model.Item {
	result := make([]*model.Item, 0, len(resp.Items))

	for _, item := range resp.Items {
		result = append(result, ItemToModel(item))
	}

	return result
}

func ItemToModel(item *inventoryV1.Item) *model.Item {
	var updatedAt *time.Time
	if item.Info.UpdatedAt != nil {
		updatedAt = lo.ToPtr(item.Info.UpdatedAt.AsTime())
	}

	return &model.Item{
		UUID:      item.Uuid,
		Name:      item.Info.Name,
		Desc:      item.Info.Desc,
		Price:     item.Info.Price,
		Category:  item.Info.Category.String(),
		CreatedAt: item.Info.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
	}
}

func ToProtoPaymentMethod(method string) paymentV1.PaymentMethod {
	switch method {
	case "CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
