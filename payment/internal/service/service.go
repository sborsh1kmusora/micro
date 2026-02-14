package service

import (
	"context"

	"github.com/sborsh1kmusora/micro/payment/internal/model"
)

type PaymentService interface {
	PayOrder(context.Context, *model.PayOrderRequest) (string, error)
}
