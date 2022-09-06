package log

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/derision-test/glock"
	"github.com/stretchr/testify/assert"
)

func TestWithAndFromContext(t *testing.T) {
	t.Run("from context without value", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, NewNilLogger(), FromContext(ctx))
	})
	t.Run("from context that has value", func(t *testing.T) {
		sink := newJSONLogger(nil)
		buffer := bytes.NewBuffer(nil)
		timestamp := time.Unix(1628115072, 0)
		sink.stream = buffer

		clock := glock.NewMockClockAt(timestamp)
		logger := newTestLogger(sink, LevelDebug, nil, clock, func() {})

		ctx := context.Background()
		ctx = WithContext(ctx, logger)

		ctxLogger := FromContext(ctx)
		assert.Same(t, logger, ctxLogger)

		ctxLogger.Info("test 1234")

		// Just make sure the message is correct, we don't care
		// about anything else since that'll tell us that we got
		// the correct logger back.
		logItem := struct {
			Message string `json:"message"`
		}{}
		assert.NoError(t, json.Unmarshal(buffer.Bytes(), &logItem))
		assert.Equal(t, "test 1234", logItem.Message)
	})
}
