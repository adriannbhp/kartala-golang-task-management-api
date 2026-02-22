package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Backup stdout
		old := os.Stdout
		defer func() { os.Stdout = old }()
		_, w, _ := os.Pipe()
		os.Stdout = w

		assert.NotPanics(t, func() {
			main()
		})
		w.Close()
	})

	t.Run("failure", func(t *testing.T) {
		// Mock exitFunc
		oldExit := exitFunc
		defer func() { exitFunc = oldExit }()
		
		exitCalled := false
		exitFunc = func(code int) {
			exitCalled = true
		}

		t.Setenv("FORCE_LOAD_SCHEMA_ERROR", "true")
		
		main()
		assert.True(t, exitCalled)
	})
}

func TestLoadSchema(t *testing.T) {
	stmts, err := LoadSchema()
	assert.NoError(t, err)
	assert.NotEmpty(t, stmts)
}
