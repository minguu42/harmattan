package database_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Begin(t *testing.T) {
	t.Run("commit", func(t *testing.T) {
		require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
			database.Users{
				{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			},
			database.Projects{},
		}))

		err := func(ctx context.Context) (err error) {
			ctx, commitOrRollback, err := c.Begin(ctx)
			if err != nil {
				return err
			}
			defer commitOrRollback(&err)

			return c.CreateProject(ctx, &domain.Project{
				ID:        "project01",
				UserID:    "user01",
				Name:      "プロジェクト1",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			})
		}(t.Context())

		require.NoError(t, err)
		tdb.Assert(t, []any{
			database.Projects{
				{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			},
		})
	})
	t.Run("rollback", func(t *testing.T) {
		require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
			database.Users{
				{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			},
			database.Projects{},
		}))

		err := func(ctx context.Context) (err error) {
			ctx, commitOrRollback, err := c.Begin(ctx)
			if err != nil {
				return err
			}
			defer commitOrRollback(&err)

			if err := c.CreateProject(ctx, &domain.Project{
				ID:        "project01",
				UserID:    "user01",
				Name:      "プロジェクト1",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			}); err != nil {
				return err
			}
			return errors.New("some error")
		}(t.Context())

		assert.Error(t, err)
		tdb.Assert(t, []any{database.Projects{}})
	})
	t.Run("nested_begin", func(t *testing.T) {
		require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
			database.Users{
				{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			},
			database.Projects{},
		}))

		err := func(ctx context.Context) (err error) {
			ctx, commitOrRollback, err := c.Begin(ctx)
			if err != nil {
				return err
			}
			defer commitOrRollback(&err)

			if err := func(ctx context.Context) (err error) {
				ctx, commitOrRollback, err := c.Begin(ctx)
				if err != nil {
					return err
				}
				defer commitOrRollback(&err)

				return c.CreateProject(ctx, &domain.Project{
					ID:        "project01",
					UserID:    "user01",
					Name:      "プロジェクト1",
					Color:     "blue",
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				})
			}(ctx); err != nil {
				return err
			}

			return c.CreateProject(ctx, &domain.Project{
				ID:        "project02",
				UserID:    "user01",
				Name:      "プロジェクト2",
				Color:     "green",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			})
		}(t.Context())

		require.NoError(t, err)
		tdb.Assert(t, []any{
			database.Projects{
				{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				{ID: "project02", UserID: "user01", Name: "プロジェクト2", Color: "green", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			},
		})
	})
	t.Run("nil_error_pointer", func(t *testing.T) {
		_, commitOrRollback, err := c.Begin(t.Context())
		require.NoError(t, err)
		assert.Panics(t, func() { commitOrRollback(nil) })
	})
}
