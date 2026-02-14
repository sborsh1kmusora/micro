package model

import "github.com/sborsh1kmusora/micro/order/internal/model"

type Order struct {
	UUID            string
	UserUUID        string
	ItemUuids       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   string
	Status          model.OrderStatus
}
