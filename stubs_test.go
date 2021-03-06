package sqlmock

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

type NullInt struct {
	Integer int
	Valid   bool
}

// Satisfy sql.Scanner interface
func (ni *NullInt) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		ni.Integer, ni.Valid = 0, false

	// FIXME int, int8, int16, int32 types are handled here but that should not
	// be necessary: only int64 is a driver.Value
	// Unfortunately, the sqlmock testsuite currently relies on that because
	// sqlmock doesn't properly limits itself internally to pure driver.Value.
	case int:
		ni.Integer, ni.Valid = v, true
	case int8:
		ni.Integer, ni.Valid = int(v), true
	case int16:
		ni.Integer, ni.Valid = int(v), true
	case int32:
		ni.Integer, ni.Valid = int(v), true

	case int64:
		const maxUint = ^uint(0)
		const minUint = 0
		const maxInt = int(maxUint >> 1)
		const minInt = -maxInt - 1

		if v > int64(maxInt) || v < int64(minInt) {
			return errors.New("value out of int range")
		}
		ni.Integer, ni.Valid = int(v), true
	case []byte:
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return err
		}
		ni.Integer, ni.Valid = n, true
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		ni.Integer, ni.Valid = n, true
	default:
		return fmt.Errorf("can't convert %T to integer", value)
	}
	return nil
}

// Satisfy sql.Valuer interface.
func (ni NullInt) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return int64(ni.Integer), nil
}

// Satisfy sql.Scanner interface
func (nt *NullTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		nt.Time, nt.Valid = time.Time{}, false
	case time.Time:
		nt.Time, nt.Valid = v, true
	default:
		return fmt.Errorf("can't convert %T to time.Time", value)
	}
	return nil
}

// Satisfy sql.Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
