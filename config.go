package log

import (
	"fmt"
	"strings"
)

type Config struct {
	LogLevel                  string            `env:"log_level" file:"log_level" default:"info"`
	LogEncoding               string            `env:"log_encoding" file:"log_encoding" default:"console"`
	LogColorize               bool              `env:"log_colorize" file:"log_colorize" default:"true"`
	LogJSONFieldNames         map[string]string `env:"log_json_field_names" file:"log_json_field_names"` // TODO - not in config
	LogInitialFields          LogFields         `env:"log_fields" file:"log_fields"`
	LogShortTime              bool              `env:"log_short_time" file:"log_short_time" default:"false"`
	LogDisplayFields          bool              `env:"log_display_fields" file:"log_display_fields" default:"true"`
	LogDisplayMultilineFields bool              `env:"log_display_multiline_fields" file:"log_display_multiline_fields" default:"false"`
	LogFieldBlacklist         []string          `env:"log_field_blacklist" file:"log_field_blacklist"`
}

var (
	ErrIllegalLevel    = fmt.Errorf("illegal log level")
	ErrIllegalEncoding = fmt.Errorf("illegal log encoding")
)

func (c *Config) PostLoad() error {
	c.LogLevel = strings.ToLower(c.LogLevel)

	if !isLegalLevel(c.LogLevel) {
		return ErrIllegalLevel
	}

	if !isLegalEncoding(c.LogEncoding) {
		return ErrIllegalEncoding
	}

	for name := range c.LogJSONFieldNames {
		if !isLegalJSONFieldName(name) {
			return fmt.Errorf("unknown JSON field name %s", name)
		}
	}

	for i, name := range c.LogFieldBlacklist {
		c.LogFieldBlacklist[i] = strings.ToLower(name)
	}

	return nil
}

func isLegalLevel(level string) bool {
	for _, whitelisted := range names {
		if level == whitelisted {
			return true
		}
	}

	return false
}

func isLegalEncoding(encoding string) bool {
	return encoding == "console" || encoding == "json"
}

func isLegalJSONFieldName(name string) bool {
	return name == "message" || name == "timestamp" || name == "level"
}
