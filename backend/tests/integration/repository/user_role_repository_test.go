package repository

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/tests/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRoleRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewUserRoleRepository(db)
	ctx := context.Background()

	t.Run("Create and ListByUsername success", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser", "roleuser@gmail.com", "1234Qwerty")

		err := repo.Create(ctx, user.Username, "user")
		require.Nil(t, err)

		roles, err := repo.ListByUsername(ctx, user.Username)
		require.Nil(t, err)
		require.Len(t, roles, 1)
		assert.Equal(t, "user", roles[0])
	})

	t.Run("Create and ListByUserID success", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser2", "roleuser2@gmail.com", "1234Qwerty")

		err := repo.Create(ctx, user.Username, "user")
		require.Nil(t, err)

		roles, err := repo.ListByUserID(ctx, user.ID)
		require.Nil(t, err)
		require.Len(t, roles, 1)
		assert.Equal(t, "user", roles[0])
	})

	t.Run("Create unique violation", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser3", "roleuser3@gmail.com", "1234Qwerty")

		err := repo.Create(ctx, user.Username, "user")
		require.Nil(t, err)

		err = repo.Create(ctx, user.Username, "user")
		require.NotNil(t, err)
		assert.ErrorIs(t, err, domain.ErrUniqueViolation)
	})

	t.Run("Create role not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser4", "roleuser4@gmail.com", "1234Qwerty")

		err := repo.Create(ctx, user.Username, "nonexistent_role")
		require.NotNil(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Delete success", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser5", "roleuser5@gmail.com", "1234Qwerty")

		err := repo.Create(ctx, user.Username, "user")
		require.Nil(t, err)

		err = repo.Delete(ctx, user.Username, "user")
		require.Nil(t, err)

		roles, err := repo.ListByUsername(ctx, user.Username)
		require.Nil(t, err)
		assert.Empty(t, roles)
	})

	t.Run("Delete not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser6", "roleuser6@gmail.com", "1234Qwerty")

		err := repo.Delete(ctx, user.Username, "user")
		require.NotNil(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("ListByUsername empty", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "roleuser7", "roleuser7@gmail.com", "1234Qwerty")

		roles, err := repo.ListByUsername(ctx, user.Username)
		require.Nil(t, err)
		assert.Empty(t, roles)
	})
}
