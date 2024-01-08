package utils

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"
)

const jsonLayout = "15:04"
const dbLayout = "15:04:05.000000"

type TimePart struct {
	null.Time
}

// UnmarshalJSON Parses the json string in the custom format
func (ct *TimePart) UnmarshalJSON(b []byte) (err error) {
	if nt, err := time.Parse(jsonLayout, strings.Trim(string(b), `"`)); err != nil {
		ct.Time = null.NewTime(time.Now(), false)
	} else {
		ct.Time = null.TimeFrom(nt)
	}
	return nil
}

// MarshalJSON writes a quoted string in the custom format
func (ct TimePart) MarshalJSON() ([]byte, error) {
	if !ct.Valid {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.String())), nil
}

// String returns the time in the custom format
func (ct *TimePart) String() string {
	if ct.IsZero() {
		return ""
	}
	return ct.ValueOrZero().Format(jsonLayout)
}

func (t TimePart) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.ValueOrZero().Format(dbLayout), nil
}

func (t *TimePart) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		t.Time = null.NewTime(time.Now(), false)
		return nil
	}

	if bv, err := driver.DefaultParameterConverter.ConvertValue(value); err == nil {
		if v, ok := bv.([]byte); ok {
			if c, e := time.Parse(dbLayout, string(v)); e != nil {
				return e
			} else {
				t.Time = null.TimeFrom(c)
			}
			return nil
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan 15:04:05 Time")
}
