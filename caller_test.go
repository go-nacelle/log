package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/efritz/glock"
	"github.com/stretchr/testify/assert"
)

func TestCallerLogger(t *testing.T)                  { testBasic(t, InitLogger) }
func TestCallerAdapter(t *testing.T)                 { testAdapter(t, InitLogger) }
func TestCallerLoggerWithFields(t *testing.T)        { testFields(t, InitLogger) }
func TestCallerLoggerWithReplayAdapter(t *testing.T) { testReplayAdapter(t, InitLogger) }
func TestCallerLoggerWithRollupAdapter(t *testing.T) { testRollupAdapter(t, InitLogger) }
func TestCallerLoggerReplay(t *testing.T)            { testReplay(t, InitLogger) }
func TestCallerLoggerRollup(t *testing.T)            { testRollup(t, InitLogger) }

func TestCallerTrimPath(t *testing.T) {
	assert.Equal(t, "", trimPath(""))
	assert.Equal(t, "/", trimPath("/"))
	assert.Equal(t, "/foo", trimPath("/foo"))
	assert.Equal(t, "foo/bar", trimPath("/foo/bar"))
	assert.Equal(t, "bar/baz", trimPath("/foo/bar/baz"))
	assert.Equal(t, "baz/bonk", trimPath("/foo/bar/baz/bonk"))
}

var (
	testFields1 = LogFields{"A": 1}
	testFields2 = LogFields{"B": 2}
	testFields3 = LogFields{"C": 3}
)

func testBasic(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		assert.Nil(t, err)

		logger.Info("X")
		logger.InfoWithFields(LogFields{"empty": false}, "Y")
		logger.Info("Z")
		logger.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	assert.Len(t, lines, 3)

	var (
		data1 = LogFields{}
		data2 = LogFields{}
		data3 = LogFields{}
	)

	assert.Nil(t, json.Unmarshal([]byte(lines[0]), &data1))
	assert.Nil(t, json.Unmarshal([]byte(lines[1]), &data2))
	assert.Nil(t, json.Unmarshal([]byte(lines[2]), &data3))

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 42

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+0), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+1), data2["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+2), data3["caller"])
}

func testReplay(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		assert.Nil(t, err)

		// Non-replayed messages are below log level - not emitted
		adapter := NewReplayAdapter(logger, LevelDebug, LevelInfo)
		adapter.Debug("X")
		adapter.InfoWithFields(LogFields{"empty": false}, "Y")
		adapter.Debug("Z")
		adapter.Replay(LevelWarning)
		adapter.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	assert.Len(t, lines, 4)

	var (
		data1 = LogFields{}
		data2 = LogFields{}
		data3 = LogFields{}
	)

	assert.Nil(t, json.Unmarshal([]byte(lines[1]), &data1))
	assert.Nil(t, json.Unmarshal([]byte(lines[2]), &data2))
	assert.Nil(t, json.Unmarshal([]byte(lines[3]), &data3))

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 78

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+0), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+1), data2["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+2), data3["caller"])
}

func testRollup(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		assert.Nil(t, err)

		clock := glock.NewMockClock()
		adapter := adaptShim(newRollupShim(logger, clock, time.Second))
		adapter.Info("A")
		adapter.Info("A")
		adapter.Info("A")
		adapter.Info("A")
		adapter.Info("A")
		clock.BlockingAdvance(time.Second)
		adapter.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	assert.Len(t, lines, 2)

	var (
		data1 = LogFields{}
		data2 = LogFields{}
	)

	assert.Nil(t, json.Unmarshal([]byte(lines[0]), &data1))
	assert.Nil(t, json.Unmarshal([]byte(lines[1]), &data2))

	// Note: this value refers to the line number containing the first instance
	// of `logger.Info("A")` in the function literal above. If code is added
	// before that line, this value must be updated.
	start := 115

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start), data2["caller"])
}

func testAdapter(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		assert.Nil(t, err)
		// Push caller stack out once for each indirection
		indirectLogger := logger.WithIndirectCaller(3)

		log3 := func(message string) { indirectLogger.Info(message) }
		log2 := func(message string) { log3(message) }
		log1 := func(message string) { log2(message) }

		log1("X")
		log1("Y")
		log1("Z")
		logger.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	assert.Len(t, lines, 3)

	var (
		data1 = LogFields{}
		data2 = LogFields{}
		data3 = LogFields{}
	)

	assert.Nil(t, json.Unmarshal([]byte(lines[0]), &data1))
	assert.Nil(t, json.Unmarshal([]byte(lines[1]), &data2))
	assert.Nil(t, json.Unmarshal([]byte(lines[2]), &data3))

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 155

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+0), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+1), data2["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+2), data3["caller"])
}

func testFields(t *testing.T, init func(*Config) (Logger, error)) {
	testBasic(t, func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return logger.WithFields(testFields1).WithFields(testFields2).WithFields(testFields3), nil
	})
}

func testReplayAdapter(t *testing.T, init func(*Config) (Logger, error)) {
	testBasic(t, func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return NewReplayAdapter(NewReplayAdapter(NewReplayAdapter(logger))), nil
	})
}

func testRollupAdapter(t *testing.T, init func(*Config) (Logger, error)) {
	testBasic(t, func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return NewRollupAdapter(NewRollupAdapter(NewRollupAdapter(logger, time.Second), time.Second), time.Second), nil
	})
}
