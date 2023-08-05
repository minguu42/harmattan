package handler

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/opepe/gen/mock"
	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/repository"
	"go.uber.org/mock/gomock"
)

func TestHandler_CreateProject(t *testing.T) {
	type args struct {
		ctx context.Context
		req *ogen.CreateProjectReq
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository, g *mock.MockIDGenerator)
		want          *ogen.Project
		wantErr       error
	}{
		{
			name: "プロジェクト1を作成する",
			args: args{
				ctx: mockCtx,
				req: &ogen.CreateProjectReq{Name: "プロジェクト1", Color: "#1A2B3C"},
			},
			prepareMockFn: func(r *mock.MockRepository, g *mock.MockIDGenerator) {
				g.EXPECT().Generate().Return("01DXF6DT000000000000000000")
				r.EXPECT().CreateProject(mockCtx, &entity.Project{
					ID:         "01DXF6DT000000000000000000",
					UserID:     "01DXF6DT000000000000000000",
					Name:       "プロジェクト1",
					Color:      "#1A2B3C",
					IsArchived: false,
					CreatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				}).Return(nil)
			},
			want: &ogen.Project{
				ID:         "01DXF6DT000000000000000000",
				Name:       "プロジェクト1",
				Color:      "#1A2B3C",
				IsArchived: false,
				CreatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:          "コンテキストからユーザを取得できない場合はエラーを返す",
			args:          args{ctx: context.Background()},
			prepareMockFn: func(r *mock.MockRepository, g *mock.MockIDGenerator) {},
			want:          nil,
			wantErr:       errUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			g := mock.NewMockIDGenerator(c)
			tt.prepareMockFn(r, g)
			h := &Handler{Repository: r, IDGenerator: g}

			got, err := h.CreateProject(tt.args.ctx, tt.args.req)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("h.CreateProject() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.CreateProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_ListProjects(t *testing.T) {
	type args struct {
		ctx    context.Context
		params ogen.ListProjectsParams
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          *ogen.Projects
		wantErr       error
	}{
		{
			name: "プロジェクト一覧を取得する",
			args: args{ctx: mockCtx},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectsByUserID(mockCtx, "01DXF6DT000000000000000000", 11, 0).
					Return([]*entity.Project{
						{
							ID:         "01DXF6DT000000000000000000",
							UserID:     "01DXF6DT000000000000000000",
							Name:       "プロジェクト1",
							Color:      "#1A2B3C",
							IsArchived: false,
							CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:         "01DXF6DT000000000000000001",
							UserID:     "01DXF6DT000000000000000000",
							Name:       "プロジェクト2",
							Color:      "#1A2B3C",
							IsArchived: true,
							CreatedAt:  time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
							UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
						},
					}, nil)
			},
			want: &ogen.Projects{
				Projects: []ogen.Project{
					{
						ID:         "01DXF6DT000000000000000000",
						Name:       "プロジェクト1",
						Color:      "#1A2B3C",
						IsArchived: false,
						CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:         "01DXF6DT000000000000000001",
						Name:       "プロジェクト2",
						Color:      "#1A2B3C",
						IsArchived: true,
						CreatedAt:  time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
						UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
					},
				},
				HasNext: false,
			},
			wantErr: nil,
		},
		{
			name:          "コンテキストからユーザを取得できない場合はエラーを返す",
			args:          args{ctx: context.Background()},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          nil,
			wantErr:       errUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			got, err := h.ListProjects(tt.args.ctx, tt.args.params)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("h.ListProjects() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.ListProjects() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_UpdateProject(t *testing.T) {
	type args struct {
		ctx    context.Context
		req    *ogen.UpdateProjectReq
		params ogen.UpdateProjectParams
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          *ogen.Project
		wantErr       error
	}{
		{
			name: "プロジェクト1を変更する",
			args: args{
				ctx: mockCtx,
				req: &ogen.UpdateProjectReq{
					Name:       ogen.OptString{Value: "新プロジェクト1", Set: true},
					Color:      ogen.OptString{Value: "#FFFFFF", Set: true},
					IsArchived: ogen.OptBool{Value: true, Set: true},
				},
				params: ogen.UpdateProjectParams{ProjectID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:         "01DXF6DT000000000000000000",
					UserID:     "01DXF6DT000000000000000000",
					Name:       "プロジェクト1",
					Color:      "#1A2B3C",
					IsArchived: false,
					CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().UpdateProject(mockCtx, &entity.Project{
					ID:         "01DXF6DT000000000000000000",
					UserID:     "01DXF6DT000000000000000000",
					Name:       "新プロジェクト1",
					Color:      "#FFFFFF",
					IsArchived: true,
					CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				}).Return(nil)
			},
			want: &ogen.Project{
				ID:         "01DXF6DT000000000000000000",
				Name:       "新プロジェクト1",
				Color:      "#FFFFFF",
				IsArchived: true,
				CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "コンテキストからユーザを取得できない場合はエラーを返す",
			args: args{
				ctx: context.Background(),
				req: &ogen.UpdateProjectReq{Name: ogen.OptString{Value: "新プロジェクト1", Set: true}},
			},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          nil,
			wantErr:       errUnauthorized,
		},
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateProjectReq{Name: ogen.OptString{Value: "新プロジェクト2", Set: true}},
				params: ogen.UpdateProjectParams{ProjectID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(nil, repository.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "指定したプロジェクトをユーザが保持していない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateProjectReq{Name: ogen.OptString{Value: "新プロジェクト2", Set: true}},
				params: ogen.UpdateProjectParams{ProjectID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000001",
					UserID:    "01DXF6DT000000000000000001",
					Name:      "プロジェクト2",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			got, err := h.UpdateProject(tt.args.ctx, tt.args.req, tt.args.params)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("h.UpdateProject() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.UpdateProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_DeleteProject(t *testing.T) {
	type args struct {
		ctx    context.Context
		params ogen.DeleteProjectParams
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          error
	}{
		{
			name: "プロジェクト1を削除する",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteProjectParams{ProjectID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
				r.EXPECT().DeleteProject(mockCtx, "01DXF6DT000000000000000000").Return(nil)
			},
			want: nil,
		},
		{
			name:          "コンテキストからユーザを取得できない場合はエラーを返す",
			args:          args{ctx: context.Background()},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          errUnauthorized,
		},
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteProjectParams{ProjectID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(nil, repository.ErrRecordNotFound)
			},
			want: errProjectNotFound,
		},
		{
			name: "指定したプロジェクトをユーザが保持していない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteProjectParams{ProjectID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000001",
					UserID:    "01DXF6DT000000000000000001",
					Name:      "プロジェクト2",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil)
			},
			want: errProjectNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			if err := h.DeleteProject(tt.args.ctx, tt.args.params); tt.want != err {
				t.Errorf("h.DeleteProject() want '%v', but '%v'", tt.want, err)
			}
		})
	}
}
