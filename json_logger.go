package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type jsonLogger struct {
	stream         io.Writer
	messageField   string
	timestampField string
	levelField     string
}

const JSONTimeFormat = "2006-01-02T15:04:05.000-0700"

func newJSONLogger(fieldNames map[string]string) *jsonLogger {
	return &jsonLogger{
		stream:         os.Stderr,
		messageField:   getField(fieldNames, "message"),
		timestampField: getField(fieldNames, "timestamp"),
		levelField:     getField(fieldNames, "level"),
	}
}

func (l *jsonLogger) Log(timestamp time.Time, level LogLevel, fields LogFields, msg string) error {
	mergedFields := fields.clone()
	mergedFields[l.messageField] = msg
	mergedFields[l.timestampField] = timestamp.Format(JSONTimeFormat)
	mergedFields[l.levelField] = level.String()

	out, err := json.Marshal(mergedFields)
	if err != nil {
		return err
	}

	fmt.Fprint(l.stream, string(out)+"\n")
	return nil
}

func getField(fieldNames map[string]string, field string) string {
	if value, ok := fieldNames[field]; ok {
		return value
	}

	return field
}
