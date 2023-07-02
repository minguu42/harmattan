package handler

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/ogen"
	mockRepository "github.com/minguu42/mtasks/pkg/repository/mock"
)

func TestHandler_GetHealth(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *ogen.GetHealthOK
	}{
		{
			name: "Version、Revisionの値はビルド時に埋め込まれるため、テスト時には空である",
			args: args{ctx: context.Background()},
			want: &ogen.GetHealthOK{
				Version:  "",
				Revision: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := mockRepository.NewMockRepository(c)
			h := &Handler{Repository: r}

			got, err := h.GetHealth(tt.args.ctx)
			if err != nil {
				t.Fatalf("h.GetHealth failed: %s", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("h.GetHealth mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
