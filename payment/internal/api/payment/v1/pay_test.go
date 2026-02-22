package v1

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/payment/internal/converter"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrder() {
	var (
		expectedUUID = gofakeit.UUID()

		req = &paymentV1.PayOrderRequest{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
	)

	s.paymentSvc.On("PayOrder", s.ctx, converter.PayOrderRequestToModel(req)).
		Return(expectedUUID, nil).
		Once()

	resp, err := s.api.PayOrder(s.ctx, req)
	s.NoError(err)
	s.Equal(expectedUUID, resp.TransactionUuid)
}
