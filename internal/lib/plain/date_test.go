package plain_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/plain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
		want  plain.Date
	}{
		{name: "normal", year: 2026, month: 7, day: 10, want: plain.NewDate(2026, 7, 10)},
		{name: "normalize_day_overflow", year: 2026, month: 2, day: 30, want: plain.NewDate(2026, 3, 2)},
		{name: "normalize_month_overflow", year: 2026, month: 13, day: 1, want: plain.NewDate(2027, 1, 1)},
		{name: "normalize_day_zero", year: 2026, month: 7, day: 0, want: plain.NewDate(2026, 6, 30)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, plain.NewDate(tt.year, tt.month, tt.day))
		})
	}
}

func TestDateOf(t *testing.T) {
	t.Parallel()

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	tests := []struct {
		name string
		t    time.Time
		want plain.Date
	}{
		{
			name: "utc",
			t:    time.Date(2026, 7, 10, 23, 59, 59, 0, time.UTC),
			want: plain.NewDate(2026, 7, 10),
		},
		{
			name: "jst_crosses_date_boundary",
			t:    time.Date(2026, 7, 10, 23, 0, 0, 0, time.UTC).In(jst),
			want: plain.NewDate(2026, 7, 11),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, plain.DateOf(tt.t))
		})
	}
}

func TestParseDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       string
		want    plain.Date
		wantErr bool
	}{
		{name: "normal", s: "2026-07-10", want: plain.NewDate(2026, 7, 10)},
		{name: "leap_day", s: "2024-02-29", want: plain.NewDate(2024, 2, 29)},
		{name: "invalid_format", s: "2026/07/10", wantErr: true},
		{name: "invalid_day", s: "2026-02-30", wantErr: true},
		{name: "datetime", s: "2026-07-10T00:00:00Z", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := plain.ParseDate(tt.s)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDate_Year(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 2026, plain.NewDate(2026, 7, 10).Year())
}

func TestDate_Month(t *testing.T) {
	t.Parallel()

	assert.Equal(t, time.Month(7), plain.NewDate(2026, 7, 10).Month())
}

func TestDate_Day(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 10, plain.NewDate(2026, 7, 10).Day())
}

func TestDate_Weekday(t *testing.T) {
	t.Parallel()

	assert.Equal(t, time.Friday, plain.NewDate(2026, 7, 10).Weekday())
}

func TestDate_AddDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		d      plain.Date
		years  int
		months int
		days   int
		want   plain.Date
	}{
		{name: "add_days", d: plain.NewDate(2026, 7, 10), days: 22, want: plain.NewDate(2026, 8, 1)},
		{name: "add_months", d: plain.NewDate(2026, 7, 10), months: 6, want: plain.NewDate(2027, 1, 10)},
		{name: "add_years", d: plain.NewDate(2026, 7, 10), years: 1, want: plain.NewDate(2027, 7, 10)},
		{name: "subtract_days", d: plain.NewDate(2026, 7, 10), days: -10, want: plain.NewDate(2026, 6, 30)},
		{name: "normalize_month_end", d: plain.NewDate(2026, 1, 31), months: 1, want: plain.NewDate(2026, 3, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.d.AddDate(tt.years, tt.months, tt.days))
		})
	}
}

func TestDate_Compare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    plain.Date
		u    plain.Date
		want int
	}{
		{name: "equal", d: plain.NewDate(2026, 7, 10), u: plain.NewDate(2026, 7, 10)},
		{name: "earlier_year", d: plain.NewDate(2025, 12, 31), u: plain.NewDate(2026, 1, 1), want: -1},
		{name: "earlier_month", d: plain.NewDate(2026, 6, 30), u: plain.NewDate(2026, 7, 1), want: -1},
		{name: "earlier_day", d: plain.NewDate(2026, 7, 9), u: plain.NewDate(2026, 7, 10), want: -1},
		{name: "later_day", d: plain.NewDate(2026, 7, 11), u: plain.NewDate(2026, 7, 10), want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.d.Compare(tt.u))
			assert.Equal(t, tt.want > 0, tt.d.After(tt.u))
			assert.Equal(t, tt.want < 0, tt.d.Before(tt.u))
		})
	}
}

func TestDate_IsZero(t *testing.T) {
	t.Parallel()

	var zero plain.Date
	assert.True(t, zero.IsZero())
	assert.False(t, plain.NewDate(2026, 7, 10).IsZero())
}

func TestDate_In(t *testing.T) {
	t.Parallel()

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	want := time.Date(2026, 7, 10, 0, 0, 0, 0, jst)
	assert.Equal(t, want, plain.NewDate(2026, 7, 10).In(jst))
}

func TestDate_Format(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "2026/07/10", plain.NewDate(2026, 7, 10).Format("2006/01/02"))
}

func TestDate_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    plain.Date
		want string
	}{
		{name: "normal", d: plain.NewDate(2026, 7, 10), want: "2026-07-10"},
		{name: "zero_padding", d: plain.NewDate(1, 1, 1), want: "0001-01-01"},
		{name: "zero_value", want: "0000-00-00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.d.String())
		})
	}
}

func TestDate_MarshalText(t *testing.T) {
	t.Parallel()

	got, err := plain.NewDate(2026, 7, 10).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, []byte("2026-07-10"), got)
}

func TestDate_UnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		want    plain.Date
		wantErr bool
	}{
		{name: "normal", data: []byte("2026-07-10"), want: plain.NewDate(2026, 7, 10)},
		{name: "invalid", data: []byte("not-a-date"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got plain.Date
			err := got.UnmarshalText(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDate_Value(t *testing.T) {
	t.Parallel()

	got, err := plain.NewDate(2026, 7, 10).Value()
	require.NoError(t, err)
	assert.Equal(t, "2026-07-10", got)
}

func TestDate_Scan(t *testing.T) {
	t.Parallel()

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	tests := []struct {
		name    string
		src     any
		want    plain.Date
		wantErr bool
	}{
		{name: "time", src: time.Date(2026, 7, 10, 0, 0, 0, 0, jst), want: plain.NewDate(2026, 7, 10)},
		{name: "bytes", src: []byte("2026-07-10"), want: plain.NewDate(2026, 7, 10)},
		{name: "string", src: "2026-07-10", want: plain.NewDate(2026, 7, 10)},
		{name: "invalid_string", src: "not-a-date", wantErr: true},
		{name: "unsupported_type", src: 20260710, wantErr: true},
		{name: "nil", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got plain.Date
			err := got.Scan(tt.src)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
