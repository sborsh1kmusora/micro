package order

import (
	"context"
	"errors"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID, paymentMethod string) (string, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", model.ErrOrderAlreadyPaidOrCancelled
	}

	transactionUUID, err := s.paymentCL.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = paymentMethod
	order.TransactionUUID = transactionUUID

	if err := s.orderRepo.Create(ctx, order.UUID, order); err != nil {
		return "", err
	}

	return transactionUUID, nil
}
