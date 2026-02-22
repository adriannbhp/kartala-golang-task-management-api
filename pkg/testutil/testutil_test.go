package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestRouter(t *testing.T) {
	r := SetupTestRouter()
	assert.NotNil(t, r)
}

func TestSetupTestDB(t *testing.T) {
	// We want to test the error paths, so we don't skip in short mode here if we can help it
	// but real DB connection needs a DB.

	t.Run("success", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping success case in short mode")
		}
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		db, err := SetupTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("error_config_load", func(t *testing.T) {
		// Unset required fields so config.LoadConfig() returns an error,
		// covering the "return nil, err" branch at testutil.go:17-18.
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("JWT_SECRET", "") // Required by config validation → triggers error

		db, err := SetupTestDB()
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("error_connection", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("DB_HOST", "invalid_host_for_test")
		
		db, err := SetupTestDB()
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("error_migrate", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("FORCE_MIGRATE_ERROR", "true")
		
		db, err := SetupTestDB()
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "forced migration error")
	})
}

func TestClearTables(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database clear test in short mode")
	}

	db, err := SetupTestDB()
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		ClearTables(db, "tasks", "users")
	})

	t.Run("empty_tables", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ClearTables(db)
		})
	})
}
