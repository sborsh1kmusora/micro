package model

import "time"

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	UUID            string
	UserUUID        string
	ItemUuids       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   string
	Status          OrderStatus
}

type Item struct {
	UUID      string
	Name      string
	Desc      string
	Price     float64
	Category  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
