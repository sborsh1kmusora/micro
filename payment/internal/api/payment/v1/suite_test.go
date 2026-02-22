package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sborsh1kmusora/micro/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	paymentSvc *mocks.PaymentService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentSvc = mocks.NewPaymentService(s.T())

	s.api = NewAPI(
		s.paymentSvc,
	)
}

func (s *APISuite) TearDownTest() {
	s.paymentSvc.AssertExpectations(s.T())
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
