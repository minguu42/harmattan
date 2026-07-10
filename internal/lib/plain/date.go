package plain

import (
	"cmp"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

// Date はタイムゾーンに関わらない日付である。
type Date struct {
	year  int
	month time.Month
	day   int
}

var (
	_ fmt.Stringer             = Date{}
	_ encoding.TextMarshaler   = Date{}
	_ encoding.TextUnmarshaler = (*Date)(nil)
	_ driver.Valuer            = Date{}
	_ sql.Scanner              = (*Date)(nil)
)

// NewDate は time.Date と同様に範囲外の値を正規化する（例：2月30日は3月2日になる）。
func NewDate(year int, month time.Month, day int) Date {
	return DateOf(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

// DateOf は t のタイムゾーンにおける日付を返す。
func DateOf(t time.Time) Date {
	year, month, day := t.Date()
	return Date{year: year, month: month, day: day}
}

// ParseDate は "2006-01-02" 形式の文字列を解釈する。
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Date{}, errtrace.Wrap(err)
	}
	return DateOf(t), nil
}

func (d Date) Year() int {
	return d.year
}

func (d Date) Month() time.Month {
	return d.month
}

func (d Date) Day() int {
	return d.day
}

func (d Date) Weekday() time.Weekday {
	return d.In(time.UTC).Weekday()
}

func (d Date) AddDate(years, months, days int) Date {
	return DateOf(d.In(time.UTC).AddDate(years, months, days))
}

func (d Date) After(u Date) bool {
	return d.Compare(u) > 0
}

func (d Date) Before(u Date) bool {
	return d.Compare(u) < 0
}

func (d Date) Compare(u Date) int {
	return cmp.Or(
		cmp.Compare(d.year, u.year),
		cmp.Compare(d.month, u.month),
		cmp.Compare(d.day, u.day),
	)
}

func (d Date) IsZero() bool {
	return d == Date{}
}

// In は d の00:00:00を表す loc の時刻を返す。
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, loc)
}

func (d Date) Format(layout string) string {
	return d.In(time.UTC).Format(layout)
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Date) UnmarshalText(data []byte) error {
	parsed, err := ParseDate(string(data))
	if err != nil {
		return errtrace.Wrap(err)
	}
	*d = parsed
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case time.Time:
		*d = DateOf(v)
		return nil
	case []byte:
		return errtrace.Wrap(d.UnmarshalText(v))
	case string:
		return errtrace.Wrap(d.UnmarshalText([]byte(v)))
	default:
		return errtrace.Wrap(fmt.Errorf("unsupported type: %T", src))
	}
}
