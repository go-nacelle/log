package log

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&LoggerSuite{})
		s.AddSuite(&CallerSuite{})
		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&GomolJSONSuite{})
		s.AddSuite(&ReplaySuite{})
		s.AddSuite(&RollupSuite{})
	})
}

//
// Mocks

type testShim struct {
	messages []*logMessage
	mutex    sync.RWMutex
}

func (ts *testShim) WithFields(fields Fields) logShim {
	return ts
}

func (ts *testShim) copy() []*logMessage {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	messages := make([]*logMessage, len(ts.messages))
	copy(messages, ts.messages)
	return messages
}

func (ts *testShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	ts.messages = append(ts.messages, &logMessage{
		level:  level,
		fields: addCaller(fields),
		format: format,
		args:   args,
	})
}

func (ts *testShim) Sync() error {
	return nil
}

//
// Log Capture

func captureStderr(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err.Error())
	}

	ch := make(chan string)
	go read(reader, ch)
	replaceStderr(writer, f)
	return <-ch
}

func read(reader io.Reader, ch chan<- string) {
	defer close(ch)

	var (
		buffer  = bytes.Buffer{}
		scanner = bufio.NewScanner(reader)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if _, err := buffer.Write([]byte(line + "\n")); err != nil {
			panic(err.Error())
		}
	}

	ch <- buffer.String()
}

func replaceStderr(writer *os.File, f func()) {
	defer writer.Close()

	temp := os.Stderr
	os.Stderr = writer
	f()
	os.Stderr = temp
}
