package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLegalLevel(t *testing.T) {
	assert.True(t, isLegalLevel("debug"))
	assert.True(t, isLegalLevel("info"))
	assert.True(t, isLegalLevel("warning"))
	assert.True(t, isLegalLevel("error"))
	assert.True(t, isLegalLevel("fatal"))
	assert.False(t, isLegalLevel("warn"))
	assert.False(t, isLegalLevel("trace"))
	assert.False(t, isLegalLevel("die"))
}

func TestIsLegalEncoding(t *testing.T) {
	assert.True(t, isLegalEncoding("json"))
	assert.True(t, isLegalEncoding("console"))
	assert.False(t, isLegalEncoding("file"))
	assert.False(t, isLegalEncoding("yaml"))
}
