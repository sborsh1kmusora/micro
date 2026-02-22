package v1

import (
	"errors"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestGetOrder() {
	tests := []struct {
		name           string
		serviceOrder   *model.Order
		serviceErr     error
		expectedType   interface{}
		expectedStatus int
	}{
		{
			name: "success",
			serviceOrder: &model.Order{
				UUID:       "order-uuid-123",
				UserUUID:   "user-uuid-123",
				ItemUuids:  []string{"item1", "item2"},
				TotalPrice: 200,
				Status:     model.OrderStatusPendingPayment,
			},
			serviceErr:   nil,
			expectedType: &orderV1.Order{},
		},
		{
			name:           "not found",
			serviceErr:     model.ErrOrderNotFound,
			expectedType:   &orderV1.NotFoundError{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal error",
			serviceErr:     errors.New("db error"),
			expectedType:   &orderV1.InternalServerError{},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderUUID := gofakeit.UUID()

			s.orderSvc.
				On("Get", s.ctx, orderUUID).
				Return(tt.serviceOrder, tt.serviceErr).
				Once()

			res, err := s.api.GetOrder(s.ctx, orderV1.GetOrderParams{
				OrderUUID: orderUUID,
			})

			s.Require().NoError(err)
			s.Require().IsType(tt.expectedType, res)

			switch v := res.(type) {

			case *orderV1.Order:
				s.Equal(tt.serviceOrder.UUID, v.UUID)
				s.Equal(tt.serviceOrder.UserUUID, v.UserUUID)
				s.Equal(tt.serviceOrder.TotalPrice, v.TotalPrice)

			case *orderV1.NotFoundError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal("Order not found", v.Message)

			case *orderV1.InternalServerError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal("Internal Server Error", v.Message)
			}
		})
	}
}
