package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("default_port_and_error", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("APP_PORT", "") // Triggers if port == "" block
		t.Setenv("DB_HOST", "invalid_host")

		err := Run()
		assert.Error(t, err)
	})

	t.Run("specific_port_and_error", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("APP_PORT", "9999")
		t.Setenv("DB_HOST", "invalid_host")

		err := Run()
		assert.Error(t, err)
	})

	t.Run("server_start_lines_coverage", func(t *testing.T) {
		// SetupApp succeeds (uses real DB from .env.test), but r.Run() is called
		// with an invalid port number so it returns an error immediately —
		// covering lines: logger.Logger.Infof(...) and return r.Run(":"+port).
		if testing.Short() {
			t.Skip("skipping server start coverage test in short mode")
		}
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("TEST_MAIN_COVERAGE", "") // ensure guard is NOT active
		t.Setenv("APP_PORT", "99999")      // invalid port → r.Run returns error

		err := Run()
		assert.Error(t, err)
	})
	t.Run("default_port_with_real_db", func(t *testing.T) {
		// SetupApp() succeeds + APP_PORT="" → covers the `if port == ""` branch.
		// TEST_MAIN_COVERAGE="true" causes Run() to return nil before r.Run() blocks.
		if testing.Short() {
			t.Skip("skipping default port coverage test in short mode")
		}
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("APP_PORT", "")              // triggers port = "8080"
		t.Setenv("TEST_MAIN_COVERAGE", "true") // exit before r.Run()

		err := Run()
		assert.NoError(t, err)
	})
}

func TestMainFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping main test in short mode")
	}

	t.Run("success", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("TEST_MAIN_COVERAGE", "true")

		assert.NotPanics(t, func() {
			main()
		})
	})

	t.Run("failure", func(t *testing.T) {
		// Mock exitFunc to prevent test exit
		oldExit := exitFunc
		defer func() { exitFunc = oldExit }()
		
		exitCalled := false
		exitFunc = func(code int) {
			exitCalled = true
		}

		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("DB_HOST", "invalid_host")
		
		main()
		assert.True(t, exitCalled)
	})
}

func TestSetupApp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping success case in short mode")
		}
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		
		r, dbCfg, err := SetupApp()
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.NotNil(t, dbCfg)
	})

	t.Run("error_config", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("JWT_SECRET", "")
		
		_, _, err := SetupApp()
		assert.Error(t, err)
	})

	t.Run("error_db", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("DB_HOST", "invalid_host")
		
		_, _, err := SetupApp()
		assert.Error(t, err)
	})
}
