package v1

import (
	"errors"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCreateOrder() {
	tests := []struct {
		name           string
		serviceUUID    string
		servicePrice   float64
		serviceErr     error
		expectedType   interface{}
		expectedStatus int
	}{
		{
			name:         "success",
			serviceUUID:  "order-uuid-123",
			servicePrice: 150.5,
			serviceErr:   nil,
			expectedType: &orderV1.CreateOrderResponse{},
		},
		{
			name:           "item not found",
			serviceErr:     model.ErrItemNotFound,
			expectedType:   &orderV1.BadRequestError{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "internal error",
			serviceErr:     errors.New("db down"),
			expectedType:   &orderV1.InternalServerError{},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			userUUID := gofakeit.UUID()
			itemUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

			s.orderSvc.
				On("Create", s.ctx, userUUID, itemUUIDs).
				Return(tt.serviceUUID, tt.servicePrice, tt.serviceErr).
				Once()

			res, err := s.api.CreateOrder(s.ctx, &orderV1.CreateOrderRequest{
				UserUUID:  userUUID,
				ItemUuids: itemUUIDs,
			})

			s.Require().NoError(err)
			s.Require().IsType(tt.expectedType, res)

			switch v := res.(type) {

			case *orderV1.CreateOrderResponse:
				s.Equal(tt.serviceUUID, v.OrderUUID)
				s.Equal(tt.servicePrice, v.TotalPrice)

			case *orderV1.BadRequestError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal(model.ErrItemNotFound.Error(), v.Message)

			case *orderV1.InternalServerError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal("Internal Server Error", v.Message)
			}
		})
	}
}
