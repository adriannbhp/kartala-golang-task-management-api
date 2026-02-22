package main

import (
	"fmt"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"io"
	"os"
	"errors"

	"ariga.io/atlas-provider-gorm/gormschema"
)

var exitFunc = os.Exit

func main() {
	stmts, err := LoadSchema()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		exitFunc(1)
	}
	io.WriteString(os.Stdout, stmts)
}

func LoadSchema() (string, error) {
	if os.Getenv("FORCE_LOAD_SCHEMA_ERROR") == "true" {
		return "", errors.New("forced load schema error")
	}
	return gormschema.New("postgres").Load(&users.User{}, &tasks.Task{})
}
