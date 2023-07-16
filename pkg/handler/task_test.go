package handler

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/mock"
	"github.com/minguu42/mtasks/gen/ogen"
	"github.com/minguu42/mtasks/pkg/entity"
	"github.com/minguu42/mtasks/pkg/repository"
)

func TestHandler_CreateTask(t *testing.T) {
	dueOn := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
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
				ctx: mockCtx,
				req: &ogen.CreateTaskReq{
					Title:    "タスク1",
					Content:  "Hello, 世界!",
					Priority: 3,
					DueOn:    ogen.OptDate{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
				},
				params: ogen.CreateTaskParams{ProjectID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:     "01DXF6DT000000000000000000",
					UserID: "01DXF6DT000000000000000000",
					Name:   "プロジェクト1",
				}, nil)
				r.EXPECT().CreateTask(mockCtx, "01DXF6DT000000000000000000", "タスク1", "Hello, 世界!", 3, &dueOn).Return(&entity.Task{
					ID:          "01DXF6DT000000000000000000",
					ProjectID:   "01DXF6DT000000000000000000",
					Title:       "タスク1",
					Content:     "Hello, 世界!",
					Priority:    3,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			want: &ogen.Task{
				ID:          "01DXF6DT000000000000000000",
				ProjectID:   "01DXF6DT000000000000000000",
				Title:       "タスク1",
				Content:     "Hello, 世界!",
				Priority:    3,
				DueOn:       ogen.OptDate{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
				CompletedAt: ogen.OptDateTime{Set: false},
				CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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
				params: ogen.CreateTaskParams{ProjectID: "01DXF6DT000000000000000001"},
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
				req:    &ogen.CreateTaskReq{Title: "タスク1"},
				params: ogen.CreateTaskParams{ProjectID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Project{
					ID:     "01DXF6DT000000000000000001",
					UserID: "01DXF6DT000000000000000001",
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
				params: ogen.ListTasksParams{ProjectID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().GetTasksByProjectID(mockCtx, "01DXF6DT000000000000000000", "-createdAt", 11, 0).
					Return([]*entity.Task{
						{
							ID:          "01DXF6DT000000000000000000",
							ProjectID:   "01DXF6DT000000000000000000",
							Title:       "タスク1",
							CompletedAt: nil,
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:          "01DXF6DT000000000000000001",
							ProjectID:   "01DXF6DT000000000000000000",
							Title:       "タスク2",
							CompletedAt: nil,
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					}, nil)
			},
			want: &ogen.Tasks{
				Tasks: []ogen.Task{
					{
						ID:          "01DXF6DT000000000000000000",
						ProjectID:   "01DXF6DT000000000000000000",
						Title:       "タスク1",
						CompletedAt: ogen.OptDateTime{Set: false},
						CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "01DXF6DT000000000000000001",
						ProjectID:   "01DXF6DT000000000000000000",
						Title:       "タスク2",
						CompletedAt: ogen.OptDateTime{Set: false},
						CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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
				params: ogen.ListTasksParams{ProjectID: "01DXF6DT000000000000000001"},
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
				params: ogen.ListTasksParams{ProjectID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000001",
					UserID:    "01DXF6DT000000000000000001",
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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
	tm1 := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
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
			tm:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			args: args{
				ctx: mockCtx,
				req: &ogen.UpdateTaskReq{
					Title:       ogen.OptString{Value: "新タスク1", Set: true},
					Content:     ogen.OptString{Value: "Goodbye", Set: true},
					Priority:    ogen.OptInt{Value: 3, Set: true},
					DueOn:       ogen.OptDate{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
					IsCompleted: ogen.OptBool{Value: true, Set: true},
				},
				params: ogen.UpdateTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:     "01DXF6DT000000000000000000",
					UserID: "01DXF6DT000000000000000000",
					Name:   "プロジェクト1",
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Task{
					ID:          "01DXF6DT000000000000000000",
					ProjectID:   "01DXF6DT000000000000000000",
					Title:       "タスク1",
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().UpdateTask(mockCtx, &entity.Task{
					ID:          "01DXF6DT000000000000000000",
					ProjectID:   "01DXF6DT000000000000000000",
					Title:       "新タスク1",
					Content:     "Goodbye",
					Priority:    3,
					DueOn:       &tm1,
					CompletedAt: &tm1,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				}).Return(nil)
			},
			want: &ogen.Task{
				ID:          "01DXF6DT000000000000000000",
				ProjectID:   "01DXF6DT000000000000000000",
				Title:       "新タスク1",
				Content:     "Goodbye",
				Priority:    3,
				DueOn:       ogen.OptDate{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
				CompletedAt: ogen.OptDateTime{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
				CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "コンテキストからユーザが取得できない場合はエラーを返す",
			args: args{
				ctx:    context.Background(),
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000000"},
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
				params: ogen.UpdateTaskParams{ProjectID: "01DXF6DT000000000000000001", TaskID: "01DXF6DT000000000000000000"},
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
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: "01DXF6DT000000000000000001", TaskID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000001",
					UserID:    "01DXF6DT000000000000000001",
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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
				params: ogen.UpdateTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, "01DXF6DT000000000000000001").Return(nil, repository.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: errTaskNotFound,
		},
		{
			name: "指定したタスクが指定したプロジェクトに含まれていない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				req:    &ogen.UpdateTaskReq{IsCompleted: ogen.OptBool{Value: true, Set: true}},
				params: ogen.UpdateTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Task{
					ID:          "01DXF6DT000000000000000001",
					ProjectID:   "01DXF6DT000000000000000001",
					Title:       "タスク2",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			want:    nil,
			wantErr: errTaskNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			h := &Handler{Repository: r}

			got, err := h.UpdateTask(tt.args.ctx, tt.args.req, tt.args.params)
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
				params: ogen.DeleteTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Task{
					ID:          "01DXF6DT000000000000000000",
					ProjectID:   "01DXF6DT000000000000000000",
					Title:       "タスク1",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().DeleteTask(mockCtx, "01DXF6DT000000000000000000").Return(nil)
			},
			want: nil,
		},
		{
			name: "コンテキストからユーザが取得できない場合はエラーを返す",
			args: args{
				ctx:    context.Background(),
				params: ogen.DeleteTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {},
			want:          errUnauthorized,
		},
		{
			name: "指定したプロジェクトが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: "01DXF6DT000000000000000001", TaskID: "01DXF6DT000000000000000000"},
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
				params: ogen.DeleteTaskParams{ProjectID: "01DXF6DT000000000000000001", TaskID: "01DXF6DT000000000000000000"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000001",
					UserID:    "01DXF6DT000000000000000001",
					Name:      "プロジェクト2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			want: errProjectNotFound,
		},
		{
			name: "指定したタスクが見つからない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, "01DXF6DT000000000000000001").Return(nil, repository.ErrRecordNotFound)
			},
			want: errTaskNotFound,
		},
		{
			name: "指定したタスクが指定したプロジェクトに含まれていない場合はエラーを返す",
			args: args{
				ctx:    mockCtx,
				params: ogen.DeleteTaskParams{ProjectID: "01DXF6DT000000000000000000", TaskID: "01DXF6DT000000000000000001"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetProjectByID(mockCtx, "01DXF6DT000000000000000000").Return(&entity.Project{
					ID:        "01DXF6DT000000000000000000",
					UserID:    "01DXF6DT000000000000000000",
					Name:      "プロジェクト1",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				r.EXPECT().GetTaskByID(mockCtx, "01DXF6DT000000000000000001").Return(&entity.Task{
					ID:          "01DXF6DT000000000000000001",
					ProjectID:   "01DXF6DT000000000000000001",
					Title:       "タスク2",
					CompletedAt: nil,
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			want: errTaskNotFound,
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
