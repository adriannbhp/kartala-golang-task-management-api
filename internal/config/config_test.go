package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("existing_env", func(t *testing.T) {
		os.Setenv("TEST_ENV", "value")
		defer os.Unsetenv("TEST_ENV")

		assert.Equal(t, "value", getEnv("TEST_ENV", "default"))
	})

	t.Run("default_env", func(t *testing.T) {
		assert.Equal(t, "default", getEnv("NON_EXISTING_ENV", "default"))
	})
}

func TestLoadConfig_Fallback(t *testing.T) {
	t.Run("jwt_secret_fallback", func(t *testing.T) {
		os.Setenv("GO_ENV", "test")
		os.Unsetenv("LOAD_ENV_IN_TEST")
		os.Unsetenv("JWT_SECRET")
		t.Setenv("API_KEY", "dummy")
		defer os.Unsetenv("GO_ENV")

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test_jwt_secret_only_for_testing_purposes_123", cfg.Secret.JwtSecret)
	})

	t.Run("api_key_fallback", func(t *testing.T) {
		os.Setenv("GO_ENV", "test")
		os.Unsetenv("LOAD_ENV_IN_TEST")
		os.Unsetenv("API_KEY")
		t.Setenv("JWT_SECRET", "dummy")
		defer os.Unsetenv("GO_ENV")

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test_api_key_123", cfg.Secret.ApiKey)
	})
}

func TestLoadConfig(t *testing.T) {
	// Original working directory to restore after tests
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	t.Run("success_basic", func(t *testing.T) {
		os.Setenv("GO_ENV", "test")
		os.Setenv("JWT_SECRET", "testsecret")
		os.Setenv("API_KEY", "testapikey")
		defer func() {
			os.Unsetenv("GO_ENV")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("API_KEY")
		}()

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "testsecret", cfg.Secret.JwtSecret)
		assert.Equal(t, "testapikey", cfg.Secret.ApiKey)
	})

	t.Run("missing_jwt_secret", func(t *testing.T) {
		tmpDir := t.TempDir()
		origWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origWd)

		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		os.Unsetenv("JWT_SECRET")
		t.Setenv("API_KEY", "dummy")

		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JWT_SECRET is required")
	})

	t.Run("missing_api_key", func(t *testing.T) {
		tmpDir := t.TempDir()
		origWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(origWd)

		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		os.Unsetenv("API_KEY")
		t.Setenv("JWT_SECRET", "dummy")

		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API_KEY is required")
	})

	t.Run("env_loading_block", func(t *testing.T) {
		// Set GO_ENV to something other than "test" to enter the loading block
		t.Setenv("GO_ENV", "not-test")

		// Case 1: No .env file anywhere (will print warning but continue)
		tmpDir := t.TempDir()
		os.Chdir(tmpDir)
		defer os.Chdir(originalWd)

		t.Setenv("JWT_SECRET", "already-set")
		t.Setenv("API_KEY", "key-set")
		os.Unsetenv("APP_PORT") // Ensure it doesn't leak from root .env
		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "already-set", cfg.Secret.JwtSecret)
		assert.Equal(t, "key-set", cfg.Secret.ApiKey)

		// Case 2: .env in current directory
		os.Unsetenv("JWT_SECRET") // Unset to allow godotenv to load it
		os.Unsetenv("API_KEY")
		err = os.WriteFile(".env", []byte("JWT_SECRET=from-current\nAPI_KEY=key-from-current\nAPP_PORT=1234"), 0644)
		assert.NoError(t, err)

		cfg, err = LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "from-current", cfg.Secret.JwtSecret)
		assert.Equal(t, "key-from-current", cfg.Secret.ApiKey)
		assert.Equal(t, "1234", cfg.App.Port)

		// Case 3: .env in parent directory
		os.Remove(".env")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("API_KEY")

		subDir := filepath.Join(tmpDir, "a", "b", "c")
		os.MkdirAll(subDir, 0755)

		err = os.WriteFile(filepath.Join(tmpDir, "a", ".env"), []byte("JWT_SECRET=from-parent\nAPI_KEY=key-from-parent"), 0644)
		assert.NoError(t, err)

		os.Chdir(subDir)
		cfg, err = LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "from-parent", cfg.Secret.JwtSecret)
		assert.Equal(t, "key-from-parent", cfg.Secret.ApiKey)

		// Case 4: .env in deep sub-parent directory (searching up)
		os.RemoveAll(filepath.Join(tmpDir, "a"))
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("API_KEY")

		deepDir := filepath.Join(tmpDir, "x", "y", "z", "w") // deep enough to trigger paths
		os.MkdirAll(deepDir, 0755)

		// Place .env in tmpDir/x which is ../../../ from deepDir
		err = os.WriteFile(filepath.Join(tmpDir, "x", ".env"), []byte("JWT_SECRET=from-deep-parent\nAPI_KEY=key-from-deep-parent"), 0644)
		assert.NoError(t, err)

		os.Chdir(deepDir)
		cfg, err = LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "from-deep-parent", cfg.Secret.JwtSecret)
		assert.Equal(t, "key-from-deep-parent", cfg.Secret.ApiKey)

		os.Chdir(originalWd)
	})
}

func TestGCSCredentials(t *testing.T) {
	t.Run("valid_base64", func(t *testing.T) {
		// Set dummy env vars to avoid LoadConfig errors
		t.Setenv("GO_ENV", "test")
		t.Setenv("JWT_SECRET", "dummy")
		t.Setenv("API_KEY", "dummy")

		// dummy-json in base64
		// {"type": "service_account"} -> eyJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCJ9
		dummyBase64 := "eyJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCJ9"
		t.Setenv("GCS_CREDENTIALS_BASE64", dummyBase64)

		_, err := LoadConfig()
		assert.NoError(t, err)

		credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		assert.NotEmpty(t, credsPath)
		assert.Contains(t, credsPath, "gcs-credentials.json")

		content, err := os.ReadFile(credsPath)
		assert.NoError(t, err)
		assert.Equal(t, `{"type": "service_account"}`, string(content))
	})

	t.Run("invalid_base64", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("JWT_SECRET", "dummy")
		t.Setenv("API_KEY", "dummy")
		t.Setenv("GCS_CREDENTIALS_BASE64", "invalid@@@")

		// Should not panic, just log warning and not set env
		os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		_, err := LoadConfig()
		assert.NoError(t, err)
		assert.Empty(t, os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	})
}
