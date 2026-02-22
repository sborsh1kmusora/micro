package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		userUUID  = gofakeit.UUID()
		item1UUID = gofakeit.UUID()
		item2UUID = gofakeit.UUID()

		itemsUUIDs = []string{item1UUID, item2UUID}

		items = []*model.Item{
			{
				UUID:  item1UUID,
				Name:  gofakeit.Name(),
				Price: 50.0,
			},
			{
				UUID:  item2UUID,
				Name:  gofakeit.Name(),
				Price: 50.0,
			},
		}

		expectedTotalPrice = 100.0
	)

	s.invClient.On("ListItems", s.ctx, itemsUUIDs).Return(items, nil).Once()
	s.orderRepo.On("Create", s.ctx, mock.AnythingOfType("string"), mock.MatchedBy(func(o *model.Order) bool {
		return o.UserUUID == userUUID &&
			o.TotalPrice == expectedTotalPrice &&
			len(o.ItemUuids) == len(itemsUUIDs)
	})).Return(nil).Once()

	orderUUID, totalPrice, err := s.svc.Create(s.ctx, userUUID, itemsUUIDs)

	s.Require().NoError(err)
	s.Require().Equal(expectedTotalPrice, totalPrice)
	s.Require().NotEmpty(orderUUID)
}

func (s *ServiceSuite) TestCreateOrderItemNotFound() {
	var (
		userUUID = gofakeit.UUID()

		itemsUUIDs = []string{gofakeit.UUID(), gofakeit.UUID()}

		invErr = status.Error(codes.NotFound, "items not found")
	)

	s.invClient.On("ListItems", s.ctx, itemsUUIDs).Return(nil, invErr).Once()

	orderUUID, totalPrice, err := s.svc.Create(s.ctx, userUUID, itemsUUIDs)

	s.Require().Error(err)
	s.Require().Equal(model.ErrItemNotFound, err)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
}

func (s *ServiceSuite) TestCreateOrderRepoError() {
	var (
		userUUID  = gofakeit.UUID()
		item1UUID = gofakeit.UUID()
		item2UUID = gofakeit.UUID()

		itemsUUIDs = []string{item1UUID, item2UUID}

		items = []*model.Item{
			{
				UUID:  item1UUID,
				Name:  gofakeit.Name(),
				Price: 50.0,
			},
			{
				UUID:  item2UUID,
				Name:  gofakeit.Name(),
				Price: 50.0,
			},
		}

		expectedTotalPrice = 100.0

		repoErr = gofakeit.Error()
	)

	s.invClient.On("ListItems", s.ctx, itemsUUIDs).Return(items, nil).Once()
	s.orderRepo.On("Create", s.ctx, mock.AnythingOfType("string"), mock.MatchedBy(func(o *model.Order) bool {
		return o.UserUUID == userUUID &&
			o.TotalPrice == expectedTotalPrice &&
			len(o.ItemUuids) == len(itemsUUIDs)
	})).Return(repoErr).Once()

	orderUUID, totalPrice, err := s.svc.Create(s.ctx, userUUID, itemsUUIDs)

	s.Require().Error(err)
	s.Require().Equal(repoErr, err)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
}
