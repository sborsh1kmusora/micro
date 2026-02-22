package v1

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sborsh1kmusora/micro/inventory/internal/converter"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestCreateSuccess() {
	var (
		protoInfo = &inventoryV1.ItemInfo{
			Name:      gofakeit.Name(),
			Price:     gofakeit.Float64(),
			CreatedAt: timestamppb.New(time.Now()),
		}

		expectedUUID = gofakeit.UUID()

		req = &inventoryV1.CreateItemRequest{
			Info: protoInfo,
		}

		itemInfo = converter.InfoProtoToModel(protoInfo)
	)

	ctx := context.Background()

	s.itemSvc.On("Create", ctx, itemInfo).Return(expectedUUID, nil)

	res, err := s.api.CreateItem(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedUUID, res.GetUuid())
}

func (s *APISuite) TestCreateValidationError() {
	var (
		protoInfo = &inventoryV1.ItemInfo{
			Name:  "",
			Price: -1,
		}

		req = &inventoryV1.CreateItemRequest{
			Info: protoInfo,
		}
	)

	res, err := s.api.CreateItem(context.Background(), req)

	s.Require().Error(err)
	s.Require().Nil(res)
}

func (s *APISuite) TestCreateServiceError() {
	var (
		protoInfo = &inventoryV1.ItemInfo{
			Name:      gofakeit.Name(),
			Price:     gofakeit.Float64(),
			CreatedAt: timestamppb.New(time.Now()),
		}

		req = &inventoryV1.CreateItemRequest{
			Info: protoInfo,
		}

		itemInfo = converter.InfoProtoToModel(protoInfo)

		svcErr = gofakeit.Error()
	)

	ctx := context.Background()

	s.itemSvc.On("Create", ctx, itemInfo).Return("", svcErr)

	res, err := s.api.CreateItem(ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
