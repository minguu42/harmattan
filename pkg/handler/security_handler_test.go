package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/opepe/gen/mock"
	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/domain/model"
	"go.uber.org/mock/gomock"
)

func TestSecurity_HandleIsAuthorized(t *testing.T) {
	type args struct {
		ctx           context.Context
		operationName string
		t             ogen.IsAuthorized
	}
	tests := []struct {
		name          string
		args          args
		prepareMockFn func(r *mock.MockRepository)
		want          *model.User
		wantErr       error
	}{
		{
			name: "ユーザを取得し、コンテキストに含める",
			args: args{
				ctx: context.Background(),
				t:   ogen.IsAuthorized{APIKey: "valid api key"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetUserByAPIKey(context.Background(), "valid api key").
					Return(&model.User{
						ID:        "01DXF6DT000000000000000000",
						Name:      "ユーザ1",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
			},
			want: &model.User{
				ID:        "01DXF6DT000000000000000000",
				Name:      "ユーザ1",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "不正なAPIキーを受け取った場合はエラーを返す",
			args: args{
				ctx: context.Background(),
				t:   ogen.IsAuthorized{APIKey: "invalid api key"},
			},
			prepareMockFn: func(r *mock.MockRepository) {
				r.EXPECT().GetUserByAPIKey(context.Background(), "invalid api key").
					Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			tt.prepareMockFn(r)
			s := &Security{Repository: r}

			ctx, err := s.HandleIsAuthorized(tt.args.ctx, tt.args.operationName, tt.args.t)
			if !errors.Is(tt.wantErr, err) {
				t.Errorf("s.HandleIsAuthorized() error want '%v', but '%v'", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				got, ok := ctx.Value(userKey{}).(*model.User)
				if !ok {
					t.Fatalf("ctx.Value(userKey{}).(*model.User) failed")
				}
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("s.HandleIsAuthorized() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
