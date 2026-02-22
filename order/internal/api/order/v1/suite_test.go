package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sborsh1kmusora/micro/order/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	orderSvc *mocks.OrderService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.orderSvc = mocks.NewOrderService(s.T())

	s.api = NewAPI(
		s.orderSvc,
	)
}

func (s *APISuite) TearDownTest() {
	s.orderSvc.AssertExpectations(s.T())
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
