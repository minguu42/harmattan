package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/minguu42/mtasks/pkg/entity"
)

var taskCmpOpt = cmpopts.IgnoreFields(entity.Task{}, "ID", "CreatedAt", "UpdatedAt")

func TestDB_CreateTask(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID int64
		title     string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Task
		wantErr error
	}{
		{
			name: "新タスクを作成する",
			args: args{
				ctx:       context.Background(),
				projectID: 1,
				title:     "新タスク",
			},
			want: &entity.Task{
				ProjectID:   1,
				Title:       "新タスク",
				CompletedAt: nil,
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

			got, err := testDB.CreateTask(tt.args.ctx, tt.args.projectID, tt.args.title)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("testDB.CreateTask error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got, taskCmpOpt); diff != "" {
				t.Errorf("testDB.CreateTask mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDB_GetTaskByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Task
		wantErr error
	}{
		{
			name: "タスク1を取得する",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &entity.Task{
				ID:          1,
				ProjectID:   1,
				Title:       "タスク1",
				CompletedAt: nil,
				CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "存在しないタスクを指定した場合はエラーを返す",
			args: args{
				ctx: context.Background(),
				id:  3,
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDB.GetTaskByID(tt.args.ctx, tt.args.id)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("testDB.GetTaskByID error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testDB.GetTaskByID mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDB_GetTasksByProjectID(t *testing.T) {
	completedAt := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	type args struct {
		ctx       context.Context
		projectID int64
		sort      string
		limit     int
		offset    int
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.Task
		wantErr error
	}{
		{
			name: "タスク一覧を取得する",
			args: args{
				ctx:       context.Background(),
				projectID: 1,
				sort:      "title",
				limit:     11,
				offset:    0,
			},
			want: []*entity.Task{
				{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          2,
					ProjectID:   1,
					Title:       "タスク2",
					CompletedAt: &completedAt,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "タスクを1つだけ取得する",
			args: args{
				ctx:       context.Background(),
				projectID: 1,
				sort:      "title",
				limit:     1,
				offset:    0,
			},
			want: []*entity.Task{
				{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDB.GetTasksByProjectID(tt.args.ctx, tt.args.projectID, tt.args.sort, tt.args.limit, tt.args.offset)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("testDB.GetTasksByProjectID error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testDB.GetTasksByProjectID mismatch (-want +got):\n%s", err)
			}
		})
	}
}

func TestDB_UpdateTask(t *testing.T) {
	completedAt := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	type args struct {
		ctx         context.Context
		id          int64
		completedAt *time.Time
		updatedAt   time.Time
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "タスク1を更新する",
			args: args{
				ctx:         context.Background(),
				id:          1,
				completedAt: &completedAt,
				updatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
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

			if err := testDB.UpdateTask(tt.args.ctx, tt.args.id, tt.args.completedAt, tt.args.updatedAt); (tt.want == nil) != (err == nil) {
				t.Errorf("testDB.UpdateTask want '%v', but '%v'", tt.want, err)
			}
		})
	}
}

func TestDB_DeleteTask(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "タスク1を削除する",
			args: args{
				ctx: context.Background(),
				id:  1,
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

			if err := testDB.DeleteTask(tt.args.ctx, tt.args.id); (tt.want == nil) != (err == nil) {
				t.Errorf("testDB.DeleteTask want '%v', but '%v'", tt.want, err)
			}
		})
	}
}
