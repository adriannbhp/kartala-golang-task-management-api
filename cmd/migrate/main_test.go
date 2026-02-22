package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSchema(t *testing.T) {
	stmts, err := LoadSchema()

	assert.NoError(t, err)
	assert.NotEmpty(t, stmts)
	assert.Contains(t, stmts, "CREATE TABLE")
	assert.Contains(t, stmts, "users")
}
