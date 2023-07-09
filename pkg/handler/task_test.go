package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/mock"
	"github.com/minguu42/mtasks/gen/ogen"
	"github.com/minguu42/mtasks/pkg/entity"
	"github.com/minguu42/mtasks/pkg/ttime"
	"gorm.io/gorm"
)

func TestHandler_CreateTask(t *testing.T) {
	type args struct {
		ctx    context.Context
		req    *ogen.CreateTaskReq
		params ogen.CreateTaskParams
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          *ogen.Task
		wantErr       error
	}{
		{
			name: "タスク1を作成する",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.CreateTaskReq{Title: "タスク1"},
				params: ogen.CreateTaskParams{ProjectID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().CreateTask(mockCtx, int64(1), "タスク1").Return(&entity.Task{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want: &ogen.Task{
				ID:          1,
				ProjectID:   1,
				Title:       "タスク1",
				CompletedAt: ogen.OptDateTime{Set: false},
				CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
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
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.CreateTaskReq{Title: "タスク1"},
				params: ogen.CreateTaskParams{ProjectID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "repository.GetProjectByIDが何らかのエラーを返した場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.CreateTaskReq{Title: "タスク1"},
				params: ogen.CreateTaskParams{ProjectID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
		{
			name: "指定したプロジェクトをユーザが保持していない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.CreateTaskReq{Title: "タスク1"},
				params: ogen.CreateTaskParams{ProjectID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(&entity.Project{
					ID:        2,
					UserID:    2,
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "repository.CreateTaskが何らかのエラーを返す場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.CreateTaskReq{Title: "タスク1"},
				params: ogen.CreateTaskParams{ProjectID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().CreateTask(mockCtx, int64(1), "タスク1").Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			got, err := h.CreateTask(tt.args.ctx, tt.args.req, tt.args.params)
			if tt.wantErr != err {
				t.Errorf("h.CreateTask() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.CreateTask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_ListTasks(t *testing.T) {
	type args struct {
		ctx    context.Context
		params ogen.ListTasksParams
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          *ogen.Tasks
		wantErr       error
	}{
		{
			name: "タスク一覧を取得する",
			args: args{
				ctx:    mockCtx,
				params: ogen.ListTasksParams{ProjectID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTasksByProjectID(mockCtx, int64(1), "-createdAt", 11, 0).
					Return([]*entity.Task{
						{
							ID:          1,
							ProjectID:   1,
							Title:       "タスク1",
							CompletedAt: nil,
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
							UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
						},
						{
							ID:          2,
							ProjectID:   1,
							Title:       "タスク2",
							CompletedAt: nil,
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
							UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
						},
					}, nil)
			},
			want: &ogen.Tasks{
				Tasks: []ogen.Task{
					{
						ID:          1,
						ProjectID:   1,
						Title:       "タスク1",
						CompletedAt: ogen.OptDateTime{Set: false},
						CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
						UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					},
					{
						ID:          2,
						ProjectID:   1,
						Title:       "タスク2",
						CompletedAt: ogen.OptDateTime{Set: false},
						CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
						UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
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
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.ListTasksParams{ProjectID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "repository.GetProjectByIDが何らかのエラーを返した場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.ListTasksParams{ProjectID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
		{
			name: "指定したプロジェクトをユーザが保持していない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.ListTasksParams{ProjectID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(&entity.Project{
					ID:        2,
					UserID:    2,
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "repository.GetTasksByProjectIDが何らかのエラーを返す場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.ListTasksParams{ProjectID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTasksByProjectID(mockCtx, int64(1), "-createdAt", 11, 0).
					Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			got, err := h.ListTasks(tt.args.ctx, tt.args.params)
			if tt.wantErr != err {
				t.Errorf("h.ListTasks() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.ListTasks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	tm1 := time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local)
	type args struct {
		ctx    context.Context
		req    *ogen.UpdateTaskReq
		params ogen.UpdateTaskParams
	}
	tests := []struct {
		name          string
		tm            time.Time
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          *ogen.Task
		wantErr       error
	}{
		{
			name: "タスク1を更新する",
			tm:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local),
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(1)).Return(&entity.Task{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().UpdateTask(mockCtx, int64(1), &tm1, time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local)).Return(nil)
			},
			want: &ogen.Task{
				ID:          1,
				ProjectID:   1,
				Title:       "タスク1",
				CompletedAt: ogen.OptDateTime{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local), Set: true},
				CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local),
			},
			wantErr: nil,
		},
		{
			name: "リクエストボディに何も含まない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Set: false}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          nil,
			wantErr:       errBadRequest,
		},
		{
			name: "コンテキストからユーザが取得できない場合はエラーを返す",
			args: args{
				ctx:    context.Background(),
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          nil,
			wantErr:       errUnauthorized,
		},
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 2, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "repository.GetProjectByIDが何らかのエラーを返した場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
		{
			name: "指定したプロジェクトをユーザが保持していない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 2, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(&entity.Project{
					ID:        2,
					UserID:    2,
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want:    nil,
			wantErr: errProjectNotFound,
		},
		{
			name: "指定したタスクが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(2)).Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: errTaskNotFound,
		},
		{
			name: "repository.GetTaskByIDが何らかのエラーを返した場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(1)).Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
		{
			name: "指定したタスクが指定したプロジェクトに含まれていない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(2)).Return(&entity.Task{
					ID:          2,
					ProjectID:   2,
					Title:       "タスク2",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want:    nil,
			wantErr: errTaskNotFound,
		},
		{
			name: "r.UpdateTaskが何らかのエラーを返す場合はエラーを返す",
			tm:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local),
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(1)).Return(&entity.Task{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().UpdateTask(mockCtx, int64(1), &tm1, time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local)).
					Return(errors.New("some error"))
			},
			want:    nil,
			wantErr: errInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ttime.FixTime(tt.args.ctx, tt.tm)

			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			got, err := h.UpdateTask(ctx, tt.args.req, tt.args.params)
			if tt.wantErr != err {
				t.Errorf("h.UpdateTask() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.UpdateTask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	type args struct {
		ctx    context.Context
		params ogen.DeleteTaskParams
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          error
	}{
		{
			name: "タスク1を削除する",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(1)).Return(&entity.Task{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().DeleteTask(mockCtx, int64(1)).Return(nil)
			},
			want: nil,
		},
		{
			name: "コンテキストからユーザが取得できない場合はエラーを返す",
			args: args{
				ctx:    context.Background(),
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          errUnauthorized,
		},
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 2, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(nil, gorm.ErrRecordNotFound)
			},
			want: errProjectNotFound,
		},
		{
			name: "repository.GetProjectByIDが何らかのエラーを返した場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(nil, errors.New("some error"))
			},
			want: errInternalServerError,
		},
		{
			name: "指定したプロジェクトをユーザが保持していない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 2, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(2)).Return(&entity.Project{
					ID:        2,
					UserID:    2,
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want: errProjectNotFound,
		},
		{
			name: "指定したタスクが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(2)).Return(nil, gorm.ErrRecordNotFound)
			},
			want: errTaskNotFound,
		},
		{
			name: "repository.GetTaskByIDが何らかのエラーを返した場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(1)).Return(nil, errors.New("some error"))
			},
			want: errInternalServerError,
		},
		{
			name: "指定したタスクが指定したプロジェクトに含まれていない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 2},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(2)).Return(&entity.Task{
					ID:          2,
					ProjectID:   2,
					Title:       "タスク2",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
			},
			want: errTaskNotFound,
		},
		{
			name: "r.DeleteTaskが何らかのエラーを返す場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: 1, TaskID: 1},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, int64(1)).Return(&entity.Project{
					ID:        1,
					UserID:    1,
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, int64(1)).Return(&entity.Task{
					ID:          1,
					ProjectID:   1,
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				}, nil)
				r.EXPECT().DeleteTask(mockCtx, int64(1)).Return(errors.New("some error"))
			},
			want: errInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			if err := h.DeleteTask(tt.args.ctx, tt.args.params); tt.want != err {
				t.Errorf("h.DeleteTask() want '%v', but '%v'", tt.want, err)
			}
		})
	}
}
