package order

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sborsh1kmusora/micro/order/internal/model"
)

func (s *service) Create(ctx context.Context, userUUID string, itemsUUIDs []string) (string, float64, error) {
	items, err := s.inventoryCl.ListItems(ctx, itemsUUIDs)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return "", 0, model.ErrItemNotFound
		}
		return "", 0, err
	}

	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.Price
	}

	newUUID := uuid.NewString()

	order := &model.Order{
		UUID:       newUUID,
		UserUUID:   userUUID,
		ItemUuids:  itemsUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}

	if err := s.orderRepo.Create(ctx, newUUID, order); err != nil {
		return "", 0, err
	}

	return newUUID, totalPrice, nil
}
