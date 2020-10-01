package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func eventually(t *testing.T, cond func() bool) bool {
	return assert.Eventually(t, cond, time.Second, 10*time.Millisecond)
}

func requireEventually(t *testing.T, cond func() bool) {
	if !eventually(t, cond) {
		t.FailNow()
	}
}
