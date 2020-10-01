package log

import (
	"bytes"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsoleLoggerrLevel(t *testing.T) {
	parsed, err := template.New("test").Parse("test: {{.message}}")
	require.Nil(t, err)

	templates := map[LogLevel]*template.Template{LevelInfo: parsed}
	logger := newConsoleLogger(templates, true)
	buffer := bytes.NewBuffer(nil)
	timestamp := time.Unix(1503939881, 0)

	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelInfo,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	assert.Equal(t, "test: test 1234\n", string(buffer.Bytes()))
}

func TestConsoleLoggerColorDisabled(t *testing.T) {
	parsed, err := template.New("test").Parse("test: {{.message}}")
	require.Nil(t, err)

	templates := map[LogLevel]*template.Template{LevelNone: parsed}
	logger := newConsoleLogger(templates, false)
	buffer := bytes.NewBuffer(nil)
	timestamp := time.Unix(1503939881, 0)

	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelInfo,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	assert.Equal(t, "test: test 1234\n", string(buffer.Bytes()))
}
