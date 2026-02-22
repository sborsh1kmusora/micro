package v1

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/inventory/internal/converter"
	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetSuccess() {
	var (
		uuid = gofakeit.UUID()

		req = &inventoryV1.GetItemRequest{
			Uuid: uuid,
		}

		info = &model.ItemInfo{
			Name:      gofakeit.Name(),
			Price:     gofakeit.Price(10, 10000),
			CreatedAt: gofakeit.Date(),
		}

		item = &model.Item{
			UUID: uuid,
			Info: info,
		}

		expectedProtoItem = converter.ItemModelToProto(item)
	)

	ctx := context.Background()

	s.itemSvc.On("Get", ctx, uuid).Return(item, nil).Once()

	res, err := s.api.GetItem(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoItem.Uuid, res.Item.Uuid)
}

func (s *APISuite) TestGetItemNotFound() {
	var (
		uuid = gofakeit.UUID()

		req = &inventoryV1.GetItemRequest{
			Uuid: uuid,
		}
	)

	ctx := context.Background()

	s.itemSvc.On("Get", ctx, uuid).Return(nil, model.ErrItemNotFound).Once()

	res, err := s.api.GetItem(ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetServiceError() {
	var (
		svcErr = gofakeit.Error()
		uuid   = gofakeit.UUID()

		req = &inventoryV1.GetItemRequest{
			Uuid: uuid,
		}
	)

	ctx := context.Background()

	s.itemSvc.On("Get", ctx, uuid).Return(nil, svcErr).Once()

	res, err := s.api.GetItem(ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
