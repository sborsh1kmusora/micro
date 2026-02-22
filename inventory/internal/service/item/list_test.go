package item

import (
	"errors"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	"github.com/sborsh1kmusora/micro/inventory/internal/repository/mocks"
)

func (s *ServiceSuite) TestList() {
	type args struct {
		uuids []string
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func()
		want      []*model.Item
		wantErr   bool
	}{
		{
			name: "success with uuids",
			args: args{
				uuids: []string{"id1", "id2"},
			},
			mockSetup: func() {
				s.itemRepo.
					On("Get", s.ctx, "id1").
					Return(&model.Item{UUID: "id1"}, nil).
					Once()

				s.itemRepo.
					On("Get", s.ctx, "id2").
					Return(&model.Item{UUID: "id2"}, nil).
					Once()
			},
			want: []*model.Item{
				{UUID: "id1"},
				{UUID: "id2"},
			},
			wantErr: false,
		},
		{
			name: "error on second get",
			args: args{
				uuids: []string{"id1", "id2"},
			},
			mockSetup: func() {
				s.itemRepo.
					On("Get", s.ctx, "id1").
					Return(&model.Item{UUID: "id1"}, nil).
					Once()

				s.itemRepo.
					On("Get", s.ctx, "id2").
					Return(nil, errors.New("repo error")).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty uuids calls list",
			args: args{
				uuids: []string{},
			},
			mockSetup: func() {
				s.itemRepo.
					On("List", s.ctx).
					Return([]*model.Item{
						{UUID: "id1"},
						{UUID: "id2"},
					}, nil).
					Once()
			},
			want: []*model.Item{
				{UUID: "id1"},
				{UUID: "id2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.itemRepo = mocks.NewItemRepository(s.T())
			s.svc = NewService(s.itemRepo)

			tt.mockSetup()

			got, err := s.svc.List(s.ctx, tt.args.uuids)

			if tt.wantErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tt.want, got)
			}

			s.itemRepo.AssertExpectations(s.T())
		})
	}
}
