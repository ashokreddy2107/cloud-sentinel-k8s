package model

import (
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserConfig(t *testing.T) {
	// Setup localized DB for this test to correctly test constraints
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
	// Need to migrate User as well for FK constraint
	DB.Exec("PRAGMA foreign_keys = ON")
	err = DB.AutoMigrate(&User{}, &UserConfig{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Helpers
	createUser := func(username string) *User {
		now := time.Now()
		u := &User{
			Username:    username,
			Enabled:     true,
			LastLoginAt: &now,
		}
		if err := DB.Create(u).Error; err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		return u
	}

	t.Run("Create UserConfig successfully", func(t *testing.T) {
		user := createUser("user1")
		config := &UserConfig{
			UserID:           user.ID,
			StorageNamespace: "ns-user1",
		}
		err := DB.Create(config).Error
		assert.NoError(t, err)
		assert.NotZero(t, config.ID)
	})

	t.Run("Fail duplicate StorageNamespace", func(t *testing.T) {
		user2 := createUser("user2")
		user3 := createUser("user3")

		config1 := &UserConfig{
			UserID:           user2.ID,
			StorageNamespace: "shared-ns",
		}
		assert.NoError(t, DB.Create(config1).Error)

		config2 := &UserConfig{
			UserID:           user3.ID,
			StorageNamespace: "shared-ns", // Duplicate
		}
		err := DB.Create(config2).Error
		assert.Error(t, err)
		// SQLite error for unique constraint
		assert.True(t, strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "unique constraint"))
	})

	t.Run("Fail duplicate UserID (via unique index on user_config user_id)", func(t *testing.T) {
		user4 := createUser("user4")

		config1 := &UserConfig{
			UserID:           user4.ID,
			StorageNamespace: "ns-4a",
		}
		assert.NoError(t, DB.Create(config1).Error)

		config2 := &UserConfig{
			UserID:           user4.ID,
			StorageNamespace: "ns-4b", // Different NS but same user
		}
		err := DB.Create(config2).Error
		assert.Error(t, err)
		// Constraint idx_user_config_user_id
		assert.True(t, strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "unique constraint"))
	})

	t.Run("GetUserConfig creates default if missing", func(t *testing.T) {
		user5 := createUser("user5")

		// Initial get - should create
		config, err := GetUserConfig(user5.ID)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotEmpty(t, config.StorageNamespace)
		assert.Equal(t, user5.ID, config.UserID)
		firstNamespace := config.StorageNamespace

		// Subsequent get - should return same
		config2, err := GetUserConfig(user5.ID)
		assert.NoError(t, err)
		assert.Equal(t, firstNamespace, config2.StorageNamespace)
	})
}
