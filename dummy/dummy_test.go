package dummy_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"project_template/dummy"

	"project_template"
	"project_template/database/dbtesting"
)

func TestDummy(t *testing.T) {

	dummy1 := dummy.Dummy{
		ID:        uuid.New(),
		Title:     "val1",
		Status:    dummy.StatusActive,
		CreatedAt: time.Now(),
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
			err := dummyRepo.Update(ctx, dummy1.ID, dummy1.Title, dummy1.Status)
			require.NoError(t, err)
		})

		t.Run("get", func(t *testing.T) {
			res, err := dummyRepo.Get(ctx, dummy1.ID)
			require.NoError(t, err)
			require.Equal(t, res.ID, dummy1.ID)
			require.Equal(t, res.Title, dummy1.Title)
			require.Equal(t, res.Status, dummy1.Status)
		})
	})
}
