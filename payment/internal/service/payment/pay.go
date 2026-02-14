package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, req *model.PayOrderRequest) (string, error) {
	transactionUUID := uuid.NewString()

	log.Printf(
		"Payment for order %s was successful with payment method %s by user %s, transaction uuid %s",
		req.OrderUUID,
		req.PaymentMethod,
		req.UserUUID,
		transactionUUID,
	)

	return transactionUUID, nil
}
