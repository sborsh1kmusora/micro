package item

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	uuid := gofakeit.UUID()

	info := &model.ItemInfo{
		Name:      gofakeit.Name(),
		Desc:      "description",
		Price:     gofakeit.Float64(),
		Category:  "CLOTHING",
		CreatedAt: time.Now(),
	}

	item := &model.Item{
		UUID: uuid,
		Info: info,
	}

	s.itemRepo.On("Get", s.ctx, uuid).Return(item, nil).Once()

	res, err := s.svc.Get(s.ctx, uuid)
	s.NoError(err)
	s.Equal(item, res)
}

func (s *ServiceSuite) TestGetItemNotFound() {
	uuid := gofakeit.UUID()
	repoError := model.ErrItemNotFound

	s.itemRepo.On("Get", s.ctx, uuid).Return(nil, repoError).Once()

	res, err := s.svc.Get(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Empty(res)
}

func (s *ServiceSuite) TestGetRepoError() {
	uuid := gofakeit.UUID()
	repoErr := gofakeit.Error()

	s.itemRepo.On("Get", s.ctx, uuid).Return(nil, repoErr).Once()

	res, err := s.svc.Get(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}
