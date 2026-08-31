package repository

import (
	"context"
	"myapp/internal/repository"
	"myapp/tests/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewCartRepository(db)
	ctx := context.Background()

	t.Run("Create and GetIDByUserID", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "user", "user@gmail.com", "1234Qwerty")

		createdCartID, err := repo.Create(ctx, user.ID)
		require.Nil(t, err)

		cartID, err := repo.GetIDByUserID(ctx, user.ID)
		require.Nil(t, err)
		assert.Equal(t, createdCartID, cartID)
	})
}
