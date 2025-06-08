// +build integration

package repository_test

import (
	"context"
	"database/sql"
	"module_6/internal/database/repository"
	"module_6/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use test database
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=notes_db_test sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	
	err = db.Ping()
	require.NoError(t, err)
	
	// Clean up before test
	_, err = db.Exec("TRUNCATE users, categories, notes CASCADE")
	require.NoError(t, err)
	
	return db
}

func TestUserRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	t.Run("Create and Get User", func(t *testing.T) {
		user := &models.User{
			Email:        "test@example.com",
			Name:         "testuser",
			PasswordHash: "hashedpassword",
		}

		// Test Create
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.False(t, user.CreatedAt.IsZero())

		// Test GetByID
		retrievedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, retrievedUser.Email)
		assert.Equal(t, user.Name, retrievedUser.Name)

		// Test GetByEmail
		retrievedUser, err = repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)

		// Test GetByName
		retrievedUser, err = repo.GetByName(ctx, user.Name)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)
	})

	t.Run("Unique Constraints", func(t *testing.T) {
		user1 := &models.User{
			Email:        "unique@example.com",
			Name:         "uniqueuser",
			PasswordHash: "hash1",
		}
		
		err := repo.Create(ctx, user1)
		assert.NoError(t, err)

		// Try to create user with same email
		user2 := &models.User{
			Email:        "unique@example.com", // Same email
			Name:         "anotheruser",
			PasswordHash: "hash2",
		}
		
		err = repo.Create(ctx, user2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrAlreadyExists)

		// Try to create user with same name
		user3 := &models.User{
			Email:        "another@example.com",
			Name:         "uniqueuser", // Same name
			PasswordHash: "hash3",
		}
		
		err = repo.Create(ctx, user3)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrAlreadyExists)
	})

	t.Run("Not Found Cases", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		assert.ErrorIs(t, err, models.ErrNotFound)

		_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.ErrorIs(t, err, models.ErrNotFound)

		_, err = repo.GetByName(ctx, "nonexistentuser")
		assert.ErrorIs(t, err, models.ErrNotFound)
	})
}