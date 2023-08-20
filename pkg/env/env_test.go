package env

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		want    *Env
		wantErr error
	}{
		{
			name: "環境変数を取得する",
			want: &Env{API: API{
				Host: "www.example.com",
				Port: 443,
			}},
			wantErr: nil,
		},
		{
			name:    "先にLoad関数を呼び出し、環境変数を読み込んでいないため、エラーを返す",
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				appEnv = *tt.want
			} else {
				appEnv = Env{}
			}

			got, err := Get()
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("Get() error want '%v', but '%v'", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
