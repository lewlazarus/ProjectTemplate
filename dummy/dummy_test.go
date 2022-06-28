package dummy_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"project_template/dummy"
	"testing"
	"time"

	"project_template"
	"project_template/database/dbtesting"
)

func TestDummy(t *testing.T) {

	dummy1 := dummy.Dummy{
		ID:        uuid.New(),
		Title:     "123",
		Status:    dummy.StatusActive,
		CreatedAt: time.Now(),
	}

	updDummy1 := dummy.Dummy{
		ID:        dummy1.ID,
		Title:     "123-UPD",
		Status:    dummy.StatusInactive,
		CreatedAt: dummy1.CreatedAt,
	}

	dbtesting.Run(t, func(ctx context.Context, t *testing.T, db project_template.DB) {

		dummyRepo := db.Dummy()

		t.Run("list", func(t *testing.T) {
			_, err := dummyRepo.List(ctx)
			require.NoError(t, err)
		})

		t.Run("create", func(t *testing.T) {
			err := dummyRepo.Create(ctx, dummy1)
			require.NoError(t, err)
		})

		t.Run("update", func(t *testing.T) {
			err := dummyRepo.Update(ctx, updDummy1.ID, updDummy1.Title, updDummy1.Status)
			require.NoError(t, err)
		})

		t.Run("get", func(t *testing.T) {
			res, err := dummyRepo.Get(ctx, updDummy1.ID)
			require.NoError(t, err)
			require.Equal(t, res.ID, updDummy1.ID)
			require.Equal(t, res.Title, updDummy1.Title)
			require.Equal(t, res.Status, updDummy1.Status)
		})
	})
}
