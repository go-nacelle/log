package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

type consoleLogger struct {
	templates map[LogLevel]*template.Template
	colorize  bool
	stream    io.Writer
}

func newConsoleLogger(templates map[LogLevel]*template.Template, colorize bool) *consoleLogger {
	return &consoleLogger{
		templates: templates,
		colorize:  colorize,
		stream:    os.Stderr,
	}
}

func (l *consoleLogger) Log(timestamp time.Time, level LogLevel, fields LogFields, msg string) error {
	if !l.colorize {
		level = LevelNone
	}

	buffer := bytes.Buffer{}
	err := l.templates[level].Execute(&buffer, map[string]interface{}{
		"timestamp":  timestamp,
		"level":      level,
		"level_name": level.String(),
		"message":    msg,
		"fields":     fields,
	})

	if err != nil {
		return err
	}

	fmt.Fprint(l.stream, buffer.String()+"\n")
	return nil
}
