package v1

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestListItemsSuccess() {
	var (
		items = []*model.Item{
			{
				UUID: "id1",
				Info: &model.ItemInfo{
					Name:  "name1",
					Price: 1.0,
				},
			},
			{
				UUID: "id2",
				Info: &model.ItemInfo{
					Name:  "name2",
					Price: 2.0,
				},
			},
		}

		uuids = []string{"id1", "id2"}

		req = &inventoryV1.ListItemsRequest{
			Uuids: uuids,
		}
	)

	ctx := context.Background()

	s.itemSvc.On("List", ctx, uuids).Return(items, nil).Once()

	res, err := s.api.ListItems(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetItems(), 2)
	s.Require().Equal("id1", res.Items[0].Uuid)
	s.Require().Equal("id2", res.Items[1].Uuid)
}

func (s *APISuite) TestListItemsNotFound() {
	req := &inventoryV1.ListItemsRequest{
		Uuids: []string{"unknown-id"},
	}

	ctx := context.Background()

	s.itemSvc.
		On("List", ctx, req.Uuids).
		Return(nil, model.ErrItemNotFound).
		Once()

	res, err := s.api.ListItems(ctx, req)

	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestListItemsServiceError() {
	req := &inventoryV1.ListItemsRequest{
		Uuids: []string{"id1"},
	}

	svcErr := gofakeit.Error()

	ctx := context.Background()

	s.itemSvc.
		On("List", ctx, req.Uuids).
		Return(nil, svcErr).
		Once()

	res, err := s.api.ListItems(ctx, req)

	s.Require().Error(err)
	s.Require().Nil(res)
	s.Require().ErrorIs(err, svcErr)
}
