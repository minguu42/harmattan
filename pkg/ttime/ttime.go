// Package ttime はテスト時に現在時刻を固定できるようにするパッケージ
package ttime

import (
	"context"
	"time"
)

// Now は渡されたコンテキストに値が含まれる場合はその値を返し、含まれない場合は現在時刻を返す
// コンテキストに固定時刻が含まれている場合はその値を返す
func Now(ctx context.Context) time.Time {
	tm, ok := ctx.Value(TimeKey{}).(time.Time)
	if !ok {
		return time.Now()
	}
	return tm
}
