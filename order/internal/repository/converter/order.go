package converter

import (
	"github.com/sborsh1kmusora/micro/order/internal/model"
	repoModel "github.com/sborsh1kmusora/micro/order/internal/repository/model"
)

func OrderToModel(o *repoModel.Order) *model.Order {
	return &model.Order{
		UUID:            o.UUID,
		UserUUID:        o.UserUUID,
		ItemUuids:       o.ItemUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   o.PaymentMethod,
		Status:          o.Status,
	}
}

func OrderModelToRepo(o *model.Order) *repoModel.Order {
	return &repoModel.Order{
		UUID:            o.UUID,
		UserUUID:        o.UserUUID,
		ItemUuids:       o.ItemUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   o.PaymentMethod,
		Status:          o.Status,
	}
}
