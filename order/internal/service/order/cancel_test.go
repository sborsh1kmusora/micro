package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	var (
		uuid = gofakeit.UUID()

		order = &model.Order{
			UUID:     uuid,
			UserUUID: gofakeit.UUID(),
			Status:   model.OrderStatusPendingPayment,
		}
	)

	s.orderRepo.On("Get", s.ctx, uuid).Return(order, nil).Once()
	s.orderRepo.On("Create", s.ctx, uuid, order).Return(nil).Once()

	err := s.svc.Cancel(s.ctx, uuid)
	s.Require().Nil(err)
	s.Require().Equal(order.Status, model.OrderStatusCancelled)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	uuid := gofakeit.UUID()

	s.orderRepo.On("Get", s.ctx, uuid).Return(nil, model.ErrOrderNotFound).Once()

	err := s.svc.Cancel(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelOrderCantBeCanceled() {
	var (
		uuid = gofakeit.UUID()

		order = &model.Order{
			UUID:     uuid,
			UserUUID: gofakeit.UUID(),
			Status:   model.OrderStatusPaid,
		}
	)

	s.orderRepo.On("Get", s.ctx, uuid).Return(order, nil).Once()

	err := s.svc.Cancel(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderCantBeCancelled)
}

func (s *ServiceSuite) TestCancelOrderRepoError() {
	var (
		uuid = gofakeit.UUID()

		order = &model.Order{
			UUID:     uuid,
			UserUUID: gofakeit.UUID(),
			Status:   model.OrderStatusPendingPayment,
		}

		repoErr = gofakeit.Error()
	)

	s.orderRepo.On("Get", s.ctx, uuid).Return(order, nil).Once()
	s.orderRepo.On("Create", s.ctx, uuid, order).Return(repoErr).Once()

	err := s.svc.Cancel(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
