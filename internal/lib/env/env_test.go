package env_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/env"
	"github.com/minguu42/harmattan/internal/lib/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("supported types", func(t *testing.T) {
		type Foo struct {
			FooField string
		}
		type Bar struct {
			BarField string
		}
		type Config struct {
			BoolField    bool
			IntField     int
			Int8Field    int8
			Int16Field   int16
			Int32Field   int32
			Int64Field   int64
			UintField    uint
			Uint8Field   uint8
			Uint16Field  uint16
			Uint32Field  uint32
			Uint64Field  uint64
			Float32Field float32
			Float64Field float64
			StringField  string
			PointerField *string
			Foo          Foo
			Bar
			TimeField     time.Time
			DurationField time.Duration
		}

		t.Setenv("BoolField", "true")
		t.Setenv("IntField", "-1")
		t.Setenv("Int8Field", "-8")
		t.Setenv("Int16Field", "-16")
		t.Setenv("Int32Field", "-32")
		t.Setenv("Int64Field", "-64")
		t.Setenv("UintField", "1")
		t.Setenv("Uint8Field", "8")
		t.Setenv("Uint16Field", "16")
		t.Setenv("Uint32Field", "32")
		t.Setenv("Uint64Field", "64")
		t.Setenv("Float32Field", "3.14")
		t.Setenv("Float64Field", "3.1415")
		t.Setenv("StringField", "Hello, World!")
		t.Setenv("PointerField", "こんにちは、世界！")
		t.Setenv("FooField", "foo")
		t.Setenv("BarField", "bar")
		t.Setenv("TimeField", "2024-02-29T12:34:56Z")
		t.Setenv("DurationField", "1h20m30s")

		want := Config{
			BoolField:     true,
			IntField:      -1,
			Int8Field:     -8,
			Int16Field:    -16,
			Int32Field:    -32,
			Int64Field:    -64,
			UintField:     1,
			Uint8Field:    8,
			Uint16Field:   16,
			Uint32Field:   32,
			Uint64Field:   64,
			Float32Field:  3.14,
			Float64Field:  3.1415,
			StringField:   "Hello, World!",
			PointerField:  ptr.Ref("こんにちは、世界！"),
			Foo:           Foo{FooField: "foo"},
			Bar:           Bar{BarField: "bar"},
			TimeField:     time.Date(2024, 2, 29, 12, 34, 56, 0, time.UTC),
			DurationField: 1*time.Hour + 20*time.Minute + 30*time.Second,
		}
		var got Config
		err := env.Load(&got)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("env tag", func(t *testing.T) {
		type Foo struct {
			Field1 string `env:"FIELD_1"`
			Field2 string `env:"-"`
			Field3 string `env:"-,"`
		}

		t.Setenv("FIELD_1", "foo")
		t.Setenv("Field1", "v1")
		t.Setenv("Field2", "v2")
		t.Setenv("Field3", "v3")
		t.Setenv("-", "hyphen")

		want := Foo{Field1: "foo", Field3: "hyphen"}
		var got Foo
		err := env.Load(&got)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("required option", func(t *testing.T) {
		type Foo struct {
			Field1 string `env:",required"`
		}

		var got Foo
		err := env.Load(&got)
		assert.Error(t, err)
	})

	t.Run("default tag", func(t *testing.T) {
		type Foo struct {
			Field1 string `default:"dv1"`
			Field2 string `default:"dv2"`
			Field3 string `env:"-" default:"dv3"`
			Field4 string `env:",required" default:"dv4"`
		}

		t.Setenv("Field1", "v1")

		want := Foo{Field1: "v1", Field2: "dv2", Field4: "dv4"}
		var got Foo
		err := env.Load(&got)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
