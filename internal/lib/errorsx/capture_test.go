package errorsx_test

import (
	"errors"
	"os"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errorsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapture(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		f    func() error
		want string
	}{
		{
			err: nil,
			f:   func() error { return nil },
		},
		{
			err:  nil,
			f:    func() error { return errors.New("new error") },
			want: "new error",
		},
		{
			err:  errors.New("some error"),
			f:    func() error { return nil },
			want: "some error",
		},
		{
			err:  errors.New("some error"),
			f:    func() error { return errors.New("new error") },
			want: "some error\nnew error",
		},
		{
			err:  errors.New("some error"),
			f:    func() error { return os.ErrClosed },
			want: "some error",
		},
	}
	for _, tt := range tests {
		err := tt.err
		errorsx.Capture(&err, tt.f)

		if tt.want == "" {
			assert.Nil(t, err)
		} else {
			require.NotNil(t, err)
			assert.Equal(t, tt.want, err.Error())
		}
	}
}
