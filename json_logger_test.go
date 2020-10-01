package log

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONLoggerLog(t *testing.T) {
	logger := newJSONLogger(nil)
	buffer := bytes.NewBuffer(nil)
	timestamp := time.Unix(1503939881, 0)
	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelFatal,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	expected := fmt.Sprintf(`{
		"level": "fatal",
		"message": "test 1234",
		"timestamp": "%s",
		"attr1": 4321
	}`, timestamp.Format(JSONTimeFormat))

	assert.JSONEq(t, expected, string(buffer.Bytes()))
}

func TestJSONLoggerCustomFieldNames(t *testing.T) {
	logger := newJSONLogger(map[string]string{
		"timestamp": "@timestamp",
		"level":     "log_level",
	})
	buffer := bytes.NewBuffer(nil)
	timestamp := time.Unix(1503939881, 0)
	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelFatal,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	expected := fmt.Sprintf(`{
		"log_level": "fatal",
		"message": "test 1234",
		"@timestamp": "%s",
		"attr1": 4321
	}`, timestamp.Format(JSONTimeFormat))

	assert.JSONEq(t, expected, string(buffer.Bytes()))
}
