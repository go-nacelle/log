package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFieldsNormalizeTimeValues(t *testing.T) {
	t1 := time.Unix(1503939881, 0)
	t2 := time.Unix(1503939891, 0)
	fields := LogFields{
		"foo":  "bar",
		"bar":  t1,
		"baz":  t2,
		"bonk": []bool{true, false, true},
	}

	// Modifies object in-place
	assert.Equal(t, fields, fields.normalizeTimeValues())

	// Non-time values remain the same
	assert.Equal(t, "bar", fields["foo"])
	assert.Equal(t, []bool{true, false, true}, fields["bonk"])

	// Times converted to ISO 8601
	tx, _ := time.Parse(JSONTimeFormat, fields["bar"].(string)) // TODO - clean this up
	assert.Equal(t, t1, tx)
	tx, _ = time.Parse(JSONTimeFormat, fields["baz"].(string)) // TODO - clean this up
	assert.Equal(t, t2, tx)
}

func TestFieldsNormalizeTimeValuesOnNilFields(t *testing.T) {
	assert.Nil(t, LogFields(nil).normalizeTimeValues())
}
