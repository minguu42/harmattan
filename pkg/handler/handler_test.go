package handler

import (
	"context"
	"testing"

	"github.com/go-faster/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/pkg/ogen"
)

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
			name: "想定されるerrが渡される",
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
			name: "想定しないerrが渡される",
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
