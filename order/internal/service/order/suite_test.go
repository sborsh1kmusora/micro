package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMock "github.com/sborsh1kmusora/micro/order/internal/client/grpc/mocks"
	repoMock "github.com/sborsh1kmusora/micro/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepo     *repoMock.OrderRepository
	invClient     *clientMock.InventoryClient
	paymentClient *clientMock.PaymentClient

	svc *service
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.orderRepo = repoMock.NewOrderRepository(s.T())

	s.invClient = clientMock.NewInventoryClient(s.T())
	s.paymentClient = clientMock.NewPaymentClient(s.T())

	s.svc = NewService(s.orderRepo, s.invClient, s.paymentClient)
}

func (s *ServiceSuite) TearDownTest() {
	s.orderRepo.AssertExpectations(s.T())
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
