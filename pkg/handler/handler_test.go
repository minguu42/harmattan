package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/ttime"
)

var mockCtx = context.WithValue(
	context.WithValue(context.Background(), ttime.TimeKey{}, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
	userKey{}, &entity.User{ID: "01DXF6DT000000000000000000", Name: "ユーザ1"})

func TestHandler_NewError(t *testing.T) {
	h := Handler{}
	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name string
		args args
		want *ogen.ErrorStatusCode
	}{
		{
			name: "ハンドラエラーが渡される",
			args: args{ctx: context.Background(), err: errUnauthorized},
			want: &ogen.ErrorStatusCode{
				StatusCode: http.StatusUnauthorized,
				Response: ogen.Error{
					Code:    http.StatusUnauthorized,
					Message: "ユーザの認証に失敗しました。もしくはユーザが認証されていません。",
				},
			},
		},
		{
			name: "ハンドラエラーでないエラーが渡される",
			args: args{ctx: context.Background(), err: errors.New("")},
			want: &ogen.ErrorStatusCode{
				StatusCode: http.StatusInternalServerError,
				Response: ogen.Error{
					Code:    http.StatusInternalServerError,
					Message: "不明なエラーが発生しました。",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.NewError(tt.args.ctx, tt.args.err)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.NewError mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
