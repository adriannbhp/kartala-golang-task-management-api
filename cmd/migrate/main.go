package main

import (
	"fmt"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := LoadSchema()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}

func LoadSchema() (string, error) {
	return gormschema.New("postgres").Load(&users.User{})
}
