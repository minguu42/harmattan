package handler

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/mock"
	"github.com/minguu42/mtasks/gen/ogen"
)

func TestHandler_GetHealth(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *ogen.GetHealthOK
		wantErr error
	}{
		{
			name: "サーバの状態を取得する",
			args: args{ctx: context.Background()},
			want: &ogen.GetHealthOK{
				Version:  "",
				Revision: "",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mock.NewMockRepository(c)
			h := &Handler{Repository: r}

			got, err := h.GetHealth(tt.args.ctx)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("h.GetHealth() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.GetHealth() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
