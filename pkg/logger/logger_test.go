package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	SetupLogger()
	assert.NotNil(t, Logger)
}
