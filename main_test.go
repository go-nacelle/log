package log

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
)

type testLogger struct {
	messages []*logMessage
	mutex    sync.RWMutex
}

func (ts *testLogger) WithFields(fields LogFields) MinimalLogger {
	return ts
}

func (ts *testLogger) copy() []*logMessage {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	messages := make([]*logMessage, len(ts.messages))
	copy(messages, ts.messages)
	return messages
}

func (ts *testLogger) LogWithFields(level LogLevel, fields LogFields, format string, args ...interface{}) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	ts.messages = append(ts.messages, &logMessage{
		level:  level,
		fields: fields,
		format: format,
		args:   args,
	})
}

func (ts *testLogger) Sync() error {
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

	buffer := bytes.Buffer{}
	scanner := bufio.NewScanner(reader)

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
