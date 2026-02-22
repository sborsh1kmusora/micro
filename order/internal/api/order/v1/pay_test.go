package v1

import (
	"errors"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestPayOrder() {
	tests := []struct {
		name           string
		serviceTxUUID  string
		serviceErr     error
		expectedType   interface{}
		expectedStatus int
	}{
		{
			name:          "success",
			serviceTxUUID: "tx-uuid-123",
			serviceErr:    nil,
			expectedType:  &orderV1.PayOrderResponse{},
		},
		{
			name:           "not found",
			serviceErr:     model.ErrOrderNotFound,
			expectedType:   &orderV1.NotFoundError{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "already paid or cancelled",
			serviceErr:     model.ErrOrderAlreadyPaidOrCancelled,
			expectedType:   &orderV1.ConflictError{},
			expectedStatus: http.StatusConflict,
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
			orderUUID := gofakeit.UUID()
			paymentMethod := orderV1.PaymentMethod("CARD")

			s.orderSvc.
				On("Pay", s.ctx, orderUUID, string(paymentMethod)).
				Return(tt.serviceTxUUID, tt.serviceErr).
				Once()

			res, err := s.api.PayOrder(
				s.ctx,
				&orderV1.PayOrderRequest{
					PaymentMethod: paymentMethod,
				},
				orderV1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			)

			s.Require().NoError(err)
			s.Require().IsType(tt.expectedType, res)

			switch v := res.(type) {

			case *orderV1.PayOrderResponse:
				s.Equal(tt.serviceTxUUID, v.TransactionUUID)

			case *orderV1.NotFoundError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal("Order not found", v.Message)

			case *orderV1.ConflictError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal("Order has already paid or cancelled", v.Message)

			case *orderV1.InternalServerError:
				s.Equal(tt.expectedStatus, v.Code)
				s.Equal("Internal Server Error", v.Message)
			}
		})
	}
}
