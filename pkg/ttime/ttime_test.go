package ttime

import (
	"context"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "コンテキストに値が含まれていないので現在時刻を返す",
			args: args{ctx: context.Background()},
			want: time.Time{},
		},
		{
			name: "コンテキストに含まれる固定時刻を返す",
			args: args{ctx: context.WithValue(context.Background(), TimeKey{}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))},
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Now(tt.args.ctx)

			// want の値がゼロ値の場合は got は現在時刻の値が含まれる
			if tt.want.IsZero() {
				if got.IsZero() {
					t.Fatalf("Now should return real value")
				}
				return
			}

			if tt.want != got {
				t.Errorf("Now want %s, but %s", tt.want.Format("2006-01-02 15:04:05"), got.Format("2006-01-02 15:04:05"))
			}
		})
	}
}
