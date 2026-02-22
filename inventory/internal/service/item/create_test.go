package item

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	itemInfo := &model.ItemInfo{
		Name:      gofakeit.Name(),
		Desc:      "description",
		Price:     gofakeit.Float64(),
		Category:  "CLOTHING",
		CreatedAt: time.Now(),
	}

	expectedUUID := gofakeit.UUID()

	s.itemRepo.On("Create", s.ctx, itemInfo).Return(expectedUUID, nil).Once()

	uuid, err := s.itemRepo.Create(s.ctx, itemInfo)
	s.Require().NoError(err)
	s.Require().Equal(expectedUUID, uuid)
}

func (s *ServiceSuite) TestCreateRepoError() {
	itemInfo := &model.ItemInfo{
		Name:      gofakeit.Name(),
		Desc:      "description",
		Price:     gofakeit.Float64(),
		Category:  "CLOTHING",
		CreatedAt: time.Now(),
	}

	repoErr := gofakeit.Error()

	s.itemRepo.On("Create", s.ctx, itemInfo).Return("", repoErr).Once()

	uuid, err := s.itemRepo.Create(s.ctx, itemInfo)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(uuid)
}
