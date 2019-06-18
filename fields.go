package log

import "time"

type LogFields map[string]interface{}

func (f LogFields) clone() LogFields {
	clone := LogFields{}
	for k, v := range f {
		clone[k] = v
	}

	return clone
}

func (f LogFields) normalizeTimeValues() LogFields {
	for key, val := range f {
		switch v := val.(type) {
		case time.Time:
			f[key] = v.Format(JSONTimeFormat)
		default:
			f[key] = v
		}
	}

	return f
}
