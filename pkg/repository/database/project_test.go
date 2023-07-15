package database

import (
	"context"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/mock"
	"github.com/minguu42/mtasks/pkg/entity"
	"github.com/minguu42/mtasks/pkg/ttime"
)

func TestDB_CreateProject(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		name   string
		color  string
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(g *mock.MockIDGenerator)
		want          *entity.Project
		wantErr       error
	}{
		{
			name: "プロジェクトを作成する",
			args: args{
				ctx:    context.WithValue(context.Background(), ttime.TimeKey{}, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
				userID: "01DXF6DT000000000000000000",
				name:   "プロジェクト",
				color:  "#1A2B3C",
			},
			prepareMockFn: func(g *mock.MockIDGenerator) {
				g.EXPECT().Generate().Return("01DXF6DT000000000000000002")
			},
			want: &entity.Project{
				ID:         "01DXF6DT000000000000000002",
				UserID:     "01DXF6DT000000000000000000",
				Name:       "プロジェクト",
				Color:      "#1A2B3C",
				IsArchived: false,
				CreatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDB.Begin(); err != nil {
				t.Fatalf("testDB.Begin failed: %s", err)
			}
			defer testDB.Rollback()

			g := mock.NewMockIDGenerator(gomock.NewController(t))
			tt.prepareMockFn(g)
			testDB.SetIDGenerator(g)

			got, err := testDB.CreateProject(tt.args.ctx, tt.args.userID, tt.args.name, tt.args.color)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("testDB.CreateProject error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testDB.CreateProject mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDB_GetProjectByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Project
		wantErr error
	}{
		{
			name: "プロジェクト1を取得する",
			args: args{
				ctx: context.Background(),
				id:  "01DXF6DT000000000000000000",
			},
			want: &entity.Project{
				ID:        "01DXF6DT000000000000000000",
				UserID:    "01DXF6DT000000000000000000",
				Name:      "プロジェクト1",
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "指定したIDのプロジェクトが存在しないので、エラーを返す",
			args: args{
				ctx: context.Background(),
				id:  "01DXF6DT000000000000000002",
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDB.GetProjectByID(tt.args.ctx, tt.args.id)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("testDB.GetProjectByID error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testDB.GetProjectByID mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDB_GetProjectsByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		sort   string
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.Project
		wantErr error
	}{
		{
			name: "プロジェクト一覧を取得する",
			args: args{
				ctx:    context.Background(),
				userID: "01DXF6DT000000000000000000",
				sort:   "name",
				limit:  11,
				offset: 0,
			},
			want: []*entity.Project{
				{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "01DXF6DT000000000000000001",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "プロジェクトを1つだけ取得する",
			args: args{
				ctx:    context.Background(),
				userID: "01DXF6DT000000000000000000",
				sort:   "name",
				limit:  1,
				offset: 0,
			},
			want: []*entity.Project{
				{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDB.GetProjectsByUserID(tt.args.ctx, tt.args.userID, tt.args.sort, tt.args.limit, tt.args.offset)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("testDB.GetProjectsByUserID error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testDB.GetProjectsByUserID mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDB_UpdateProject(t *testing.T) {
	type args struct {
		ctx       context.Context
		id        string
		name      string
		updatedAt time.Time
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "プロジェクト1を更新する",
			args: args{
				ctx:       context.Background(),
				id:        "01DXF6DT000000000000000000",
				name:      "新プロジェクト1",
				updatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDB.Begin(); err != nil {
				t.Fatalf("testDB.Begin failed: %s", err)
			}
			defer testDB.Rollback()

			if err := testDB.UpdateProject(tt.args.ctx, tt.args.id, tt.args.name, tt.args.updatedAt); (tt.want == nil) != (err == nil) {
				t.Errorf("testDB.UpdateProject want '%v', but '%v'", tt.want, err)
			}
		})
	}
}

func TestDB_DeleteProject(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "プロジェクト1を削除する",
			args: args{
				ctx: context.Background(),
				id:  "01DXF6DT000000000000000000",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDB.Begin(); err != nil {
				t.Fatalf("testDB.Begin failed: %s", err)
			}
			defer testDB.Rollback()

			if err := testDB.DeleteProject(tt.args.ctx, tt.args.id); (tt.want == nil) != (err == nil) {
				t.Errorf("testDB.DeleteProject want '%v', but '%v'", tt.want, err)
			}
		})
	}
}
