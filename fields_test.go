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

	assertTime := func(value interface{}, expected interface{}) {
		actual, _ := time.Parse(JSONTimeFormat, value.(string))
		assert.Equal(t, expected, actual)
	}

	// Times converted to ISO 8601
	assertTime(fields["bar"], t1)
	assertTime(fields["baz"], t2)
}

func TestFieldsNormalizeTimeValuesOnNilFields(t *testing.T) {
	assert.Nil(t, LogFields(nil).normalizeTimeValues())
}
