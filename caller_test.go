package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/derision-test/glock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallerLogger(t *testing.T)                 { testBasic(t, InitLogger) }
func TestCallerLoggerWithFields(t *testing.T)       { testFields(t, InitLogger) }
func TestCallerLoggerWithReplayLogger(t *testing.T) { testReplayLogger(t, InitLogger) }
func TestCallerLoggerWithRollupLogger(t *testing.T) { testRollupLogger(t, InitLogger) }
func TestCallerLoggerReplay(t *testing.T)           { testReplay(t, InitLogger) }
func TestCallerLoggerRollup(t *testing.T)           { testRollup(t, InitLogger) }
func TestCallerIndirect(t *testing.T)               { testIndirect(t, InitLogger) }

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
		require.Nil(t, err)

		logger.Info("X")
		logger.InfoWithFields(LogFields{"empty": false}, "Y")
		logger.Info("Z")
		logger.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	require.Len(t, lines, 3)

	data1 := LogFields{}
	data2 := LogFields{}
	data3 := LogFields{}
	require.Nil(t, json.Unmarshal([]byte(lines[0]), &data1))
	require.Nil(t, json.Unmarshal([]byte(lines[1]), &data2))
	require.Nil(t, json.Unmarshal([]byte(lines[2]), &data3))

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 43

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+0), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+1), data2["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+2), data3["caller"])
}

func testReplay(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		require.Nil(t, err)

		// Non-replayed messages are below log level - not emitted
		replayLogger := NewReplayLogger(logger, LevelDebug, LevelInfo)
		replayLogger.Debug("X")
		replayLogger.InfoWithFields(LogFields{"empty": false}, "Y")
		replayLogger.Debug("Z")
		replayLogger.Replay(LevelWarning)
		replayLogger.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	require.Len(t, lines, 4)

	data1 := LogFields{}
	data2 := LogFields{}
	data3 := LogFields{}
	require.Nil(t, json.Unmarshal([]byte(lines[1]), &data1))
	require.Nil(t, json.Unmarshal([]byte(lines[2]), &data2))
	require.Nil(t, json.Unmarshal([]byte(lines[3]), &data3))

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 76

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+0), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+1), data2["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start+2), data3["caller"])
}

func testRollup(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		require.Nil(t, err)

		clock := glock.NewMockClock()
		rollupLogger := FromMinimalLogger(newRollupLogger(logger, clock, time.Second))
		rollupLogger.Info("A")
		rollupLogger.Info("A")
		rollupLogger.Info("A")
		rollupLogger.Info("A")
		rollupLogger.Info("A")
		clock.BlockingAdvance(time.Second)
		rollupLogger.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	require.Len(t, lines, 2)

	data1 := LogFields{}
	data2 := LogFields{}
	require.Nil(t, json.Unmarshal([]byte(lines[0]), &data1))
	require.Nil(t, json.Unmarshal([]byte(lines[1]), &data2))

	// Note: this value refers to the line number containing the first instance
	// of `logger.Info("A")` in the function literal above. If code is added
	// before that line, this value must be updated.
	start := 110

	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start), data1["caller"])
	assert.Equal(t, fmt.Sprintf("log/caller_test.go:%d", start), data2["caller"])
}

func testIndirect(t *testing.T, init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		require.Nil(t, err)
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
	require.Len(t, lines, 3)

	data1 := LogFields{}
	data2 := LogFields{}
	data3 := LogFields{}
	require.Nil(t, json.Unmarshal([]byte(lines[0]), &data1))
	require.Nil(t, json.Unmarshal([]byte(lines[1]), &data2))
	require.Nil(t, json.Unmarshal([]byte(lines[2]), &data3))

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 147

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

func testReplayLogger(t *testing.T, init func(*Config) (Logger, error)) {
	testBasic(t, func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return NewReplayLogger(NewReplayLogger(NewReplayLogger(logger))), nil
	})
}

func testRollupLogger(t *testing.T, init func(*Config) (Logger, error)) {
	testBasic(t, func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return NewRollupLogger(NewRollupLogger(NewRollupLogger(logger, time.Second), time.Second), time.Second), nil
	})
}
