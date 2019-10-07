package log

import (
	"bytes"
	"text/template"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConsoleLoggerSuite struct{}

func (s *ConsoleLoggerSuite) TestLevel(t sweet.T) {
	parsed, err := template.New("test").Parse("test: {{.message}}")
	Expect(err).To(BeNil())

	var (
		templates = map[LogLevel]*template.Template{LevelInfo: parsed}
		logger    = newConsoleLogger(templates, true)
		buffer    = bytes.NewBuffer(nil)
		timestamp = time.Unix(1503939881, 0)
	)

	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelInfo,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(Equal("test: test 1234\n"))
}

func (s *ConsoleLoggerSuite) TestColorDisabled(t sweet.T) {
	parsed, err := template.New("test").Parse("test: {{.message}}")
	Expect(err).To(BeNil())

	var (
		templates = map[LogLevel]*template.Template{LevelNone: parsed}
		logger    = newConsoleLogger(templates, false)
		buffer    = bytes.NewBuffer(nil)
		timestamp = time.Unix(1503939881, 0)
	)

	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelInfo,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(Equal("test: test 1234\n"))
}
