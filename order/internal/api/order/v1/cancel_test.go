package v1

import (
	"errors"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCancelOrder() {
	tests := []struct {
		name           string
		serviceError   error
		expectedType   interface{}
		expectedStatus int
	}{
		{
			name:         "success",
			serviceError: nil,
			expectedType: &orderV1.CancelOrderNoContent{},
		},
		{
			name:           "not found",
			serviceError:   model.ErrOrderNotFound,
			expectedType:   &orderV1.NotFoundError{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "conflict",
			serviceError:   model.ErrOrderCantBeCancelled,
			expectedType:   &orderV1.ConflictError{},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "internal",
			serviceError:   errors.New("unexpected"),
			expectedType:   &orderV1.InternalServerError{},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderUUID := gofakeit.UUID()

			s.orderSvc.
				On("Cancel", s.ctx, orderUUID).
				Return(tt.serviceError).
				Once()

			res, err := s.api.CancelOrder(s.ctx, orderV1.CancelOrderParams{
				OrderUUID: orderUUID,
			})

			s.Require().NoError(err)
			s.Require().IsType(tt.expectedType, res)

			if tt.serviceError != nil {
				switch v := res.(type) {
				case *orderV1.NotFoundError:
					s.Equal(tt.expectedStatus, v.Code)
				case *orderV1.ConflictError:
					s.Equal(tt.expectedStatus, v.Code)
				case *orderV1.InternalServerError:
					s.Equal(tt.expectedStatus, v.Code)
				}
			}
		})
	}
}
