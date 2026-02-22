package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		uuid = gofakeit.UUID()

		order = &model.Order{
			UUID: uuid,
		}
	)

	s.orderRepo.On("Get", s.ctx, uuid).Return(order, nil).Once()

	res, err := s.svc.Get(s.ctx, uuid)
	s.Require().NoError(err)
	s.Require().Equal(order, res)
	s.Require().Equal(uuid, res.UUID)
}

func (s *ServiceSuite) TestGetOrderNotFound() {
	uuid := gofakeit.UUID()

	s.orderRepo.On("Get", s.ctx, uuid).Return(nil, model.ErrOrderNotFound).Once()

	res, err := s.svc.Get(s.ctx, uuid)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		uuid    = gofakeit.UUID()
		repoErr = gofakeit.Error()
	)

	s.orderRepo.On("Get", s.ctx, uuid).Return(nil, repoErr).Once()

	res, err := s.svc.Get(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
