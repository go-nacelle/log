package log

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type JSONLoggerSuite struct{}

func (s *JSONLoggerSuite) TestLog(t sweet.T) {
	var (
		logger    = newJSONLogger(nil)
		buffer    = bytes.NewBuffer(nil)
		timestamp = time.Unix(1503939881, 0)
	)

	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelFatal,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(fmt.Sprintf(`{
		"level": "fatal",
		"message": "test 1234",
		"timestamp": "%s",
		"attr1": 4321
	}`, timestamp.Format(JSONTimeFormat))))
}

func (s *JSONLoggerSuite) TestCustomFieldNames(t sweet.T) {
	var (
		logger = newJSONLogger(map[string]string{
			"timestamp": "@timestamp",
			"level":     "log_level",
		})
		buffer    = bytes.NewBuffer(nil)
		timestamp = time.Unix(1503939881, 0)
	)

	logger.stream = buffer

	logger.Log(
		timestamp,
		LevelFatal,
		LogFields{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(fmt.Sprintf(`{
		"log_level": "fatal",
		"message": "test 1234",
		"@timestamp": "%s",
		"attr1": 4321
	}`, timestamp.Format(JSONTimeFormat))))
}
