package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		paymentMethod = "PAYMENT_METHOD_CARD"

		order = &model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPendingPayment,
		}

		expectedTransactionUUID = gofakeit.UUID()
	)

	s.orderRepo.On("Get", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder",
		s.ctx,
		orderUUID,
		userUUID,
		paymentMethod,
	).Return(expectedTransactionUUID, nil).Once()
	s.orderRepo.On("Create", s.ctx, orderUUID, order).Return(nil).Once()

	transactionUUID, err := s.svc.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(expectedTransactionUUID, transactionUUID)
	s.Require().Equal(order.Status, model.OrderStatusPaid)
	s.Require().Equal(order.PaymentMethod, paymentMethod)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	var (
		orderUUID = gofakeit.UUID()

		paymentMethod = "PAYMENT_METHOD_CARD"
	)

	s.orderRepo.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound).Once()

	transactionUUID, err := s.svc.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderAlreadyPaidOrCancelled() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		paymentMethod = "PAYMENT_METHOD_CARD"

		order = &model.Order{
			UUID:     orderUUID,
			UserUUID: userUUID,
			Status:   model.OrderStatusPaid,
		}
	)

	s.orderRepo.On("Get", s.ctx, orderUUID).Return(order, nil).Once()

	transactionUUID, err := s.svc.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaidOrCancelled)
	s.Require().Empty(transactionUUID)
}
