package item

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sborsh1kmusora/micro/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	itemRepo *mocks.ItemRepository

	svc *service
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.itemRepo = mocks.NewItemRepository(s.T())

	s.svc = NewService(s.itemRepo)
}

func (s *ServiceSuite) TearDownTest() {
	s.itemRepo.AssertExpectations(s.T())
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
