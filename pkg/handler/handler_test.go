package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/ogen"
	"github.com/minguu42/mtasks/pkg/entity"
	"github.com/minguu42/mtasks/pkg/ttime"
)

var mockCtx = context.WithValue(
	context.WithValue(context.Background(), ttime.TimeKey{}, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
	userKey{}, &entity.User{ID: 1, Name: "ユーザ1"})

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
			name: "特定のハンドラエラーが渡される",
			args: args{ctx: context.Background(), err: errBadRequest},
			want: &ogen.ErrorStatusCode{
				StatusCode: 400,
				Response: ogen.Error{
					Message: "入力に誤りがあります。入力をご確認ください。",
					Debug:   "there is an input error",
				},
			},
		},
		{
			name: "ハンドラエラーでないエラーが渡される",
			args: args{ctx: context.Background(), err: errors.New("")},
			want: &ogen.ErrorStatusCode{
				StatusCode: 500,
				Response: ogen.Error{
					Message: "不明なエラーが発生しました。もう一度お試しください。",
					Debug:   "some error occurred on the server",
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
