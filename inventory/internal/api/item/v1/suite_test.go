package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sborsh1kmusora/micro/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	itemSvc *mocks.ItemService

	api *api
}

func (s *APISuite) SetupTest() {
	s.itemSvc = mocks.NewItemService(s.T())

	s.api = NewAPI(
		s.itemSvc,
	)
}

func (s *APISuite) TearDownTest() {
	s.itemSvc.AssertExpectations(s.T())
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
