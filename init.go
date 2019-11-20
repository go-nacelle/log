package log

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/mgutz/ansi"
)

func InitLogger(c *Config) (Logger, error) {
	baseLogger, err := initBaseLogger(c)
	if err != nil {
		return nil, err
	}

	return newBaseShim(baseLogger, parseLogLevel(c.LogLevel), c.LogInitialFields), nil
}

func initBaseLogger(c *Config) (baseLogger, error) {
	if c.LogEncoding == "json" {
		return newJSONLogger(c.LogJSONFieldNames), nil
	}

	tpl, err := newConsoleTemplate(
		c.LogShortTime,
		c.LogDisplayFields,
		c.LogDisplayMultilineFields,
		c.LogFieldBlacklist,
	)

	if err != nil {
		return nil, err
	}

	return newConsoleLogger(tpl, c.LogColorize), nil
}

func newConsoleTemplate(
	shortTime bool,
	displayFields bool,
	displayMultilineFields bool,
	blacklist []string,
) (map[LogLevel]*template.Template, error) {
	var (
		fieldPrefix  = " "
		fieldPadding = ""
		fieldSuffix  = ""
	)

	if displayMultilineFields {
		fieldPrefix = "\n    "
		fieldPadding = " "
		fieldSuffix = "\n"
	}

	fieldsTemplate := fmt.Sprintf(
		""+
			`{{if .fields}}`+
			`{{range $key, $val := .fields}}`+
			`{{if shouldDisplayAttr $key}}`+
			`%s{{$key}}%s=%s{{$val}}`+
			`{{end}}`+
			`{{end}}`+
			`%s`+
			`{{end}}`,
		fieldPrefix,
		fieldPadding,
		fieldPadding,
		fieldSuffix,
	)

	timeFormat := "2006/01/02 15:04:05.000"
	if shortTime {
		timeFormat = "15:04:05"
	}

	text :=
		"" +
			`{{color}}` +
			`[{{uppercase .levelName | printf "%1.1s"}}] ` +
			fmt.Sprintf(`[{{.timestamp.Format "%s"}}] {{.message}}`, timeFormat) +
			`{{reset}}`

	if displayFields {
		text += fieldsTemplate
	}

	colors := map[LogLevel]string{
		LevelFatal:   ansi.ColorCode("red+b"),
		LevelError:   ansi.ColorCode("red"),
		LevelWarning: ansi.ColorCode("yellow"),
		LevelInfo:    ansi.ColorCode("green"),
		LevelDebug:   ansi.ColorCode("cyan"),
		LevelNone:    "",
	}

	templates := map[LogLevel]*template.Template{}
	for level, color := range colors {
		functions := template.FuncMap{
			"color":             color,
			"reset":             ansi.ColorCode("reset"),
			"uppercase":         strings.ToUpper,
			"shouldDisplayAttr": shouldDisplayAttr(blacklist),
		}

		if color == "" {
			functions["reset"] = ""
		}

		parsed, err := template.New(level.String()).Funcs(functions).Parse(text)
		if err != nil {
			return nil, err
		}

		templates[level] = parsed
	}

	return templates, nil
}

func shouldDisplayAttr(blacklist []string) func(string) bool {
	return func(attr string) bool {
		for _, cmp := range blacklist {
			if cmp == attr {
				return false
			}
		}

		return true
	}
}
