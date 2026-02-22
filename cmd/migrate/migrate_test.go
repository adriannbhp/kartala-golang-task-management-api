package main

import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"testing"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/stretchr/testify/assert"
)

func TestGormSchemaLoad(t *testing.T) {
	// Smoke test to ensure the schema loader doesn't panic and returns something
	stmts, err := gormschema.New("postgres").Load(&users.User{})
	assert.NoError(t, err)
	assert.NotEmpty(t, stmts)
}
