// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// SmallIntBool represents a DB2 SMALLINT column used for boolean values (0 or 1).
// Since DB2 has no native BOOLEAN type, we use SMALLINT with this custom type
// that implements sql.Scanner and driver.Valuer for seamless Bun integration.
//
// Usage in struct:
//
//	type Employee struct {
//	    Active SmallIntBool `bun:",default:1"`
//	}
type SmallIntBool int32

// Scan implements sql.Scanner interface to convert DB2 SMALLINT values to SmallIntBool.
// Accepts int32, int64, int, float64, string, and byte slice representations.
func (s *SmallIntBool) Scan(val interface{}) error {
	if val == nil {
		*s = 0
		return nil
	}

	switch v := val.(type) {
	case int32:
		*s = SmallIntBool(v)
	case int64:
		*s = SmallIntBool(v)
	case int:
		*s = SmallIntBool(v)
	case float64:
		*s = SmallIntBool(int32(v))
	case string:
		// Handle string representations ("0", "1", "true", "false", etc.)
		val, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			// Try parsing as boolean string
			switch v {
			case "true", "True", "TRUE", "yes", "Yes", "YES":
				*s = 1
			case "false", "False", "FALSE", "no", "No", "NO":
				*s = 0
			default:
				return fmt.Errorf("cannot scan %q into SmallIntBool", v)
			}
		} else {
			*s = SmallIntBool(val)
		}
	case []byte:
		// Handle byte slice representations
		return s.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into SmallIntBool", val)
	}
	return nil
}

// Value implements driver.Valuer interface to convert SmallIntBool back to int64.
// driver.Value only allows int64 (not int32), per database/sql/driver contract.
func (s SmallIntBool) Value() (driver.Value, error) {
	return int64(s), nil
}

// Bool converts SmallIntBool to native Go bool.
func (s SmallIntBool) Bool() bool {
	return s != 0
}

// NullSmallIntBool represents a nullable DB2 SMALLINT boolean value.
// This is equivalent to sql.NullInt32 but specifically typed for boolean semantics.
//
// Usage in struct:
//
//	type Product struct {
//	    Available NullSmallIntBool `bun:""`
//	}
type NullSmallIntBool struct {
	SmallIntBool SmallIntBool
	Valid        bool
}

// Scan implements sql.Scanner interface for NullSmallIntBool.
func (ns *NullSmallIntBool) Scan(val interface{}) error {
	if val == nil {
		ns.SmallIntBool = 0
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return ns.SmallIntBool.Scan(val)
}

// Value implements driver.Valuer interface for NullSmallIntBool.
func (ns NullSmallIntBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.SmallIntBool.Value()
}

// Bool converts NullSmallIntBool to *bool (nil if not valid).
func (ns NullSmallIntBool) Bool() *bool {
	if !ns.Valid {
		return nil
	}
	b := ns.SmallIntBool.Bool()
	return &b
}

// SmallInt represents a DB2 SMALLINT column as int32.
// Use this for non-boolean integer fields that map to SMALLINT.
//
// Usage in struct:
//
//	type DataType struct {
//	    Quantity SmallInt `bun:""`
//	}
type SmallInt int32

// Scan implements sql.Scanner interface for SmallInt.
func (si *SmallInt) Scan(val interface{}) error {
	if val == nil {
		*si = 0
		return nil
	}

	switch v := val.(type) {
	case int32:
		*si = SmallInt(v)
	case int64:
		*si = SmallInt(v)
	case int:
		*si = SmallInt(v)
	case float64:
		*si = SmallInt(int32(v))
	case string:
		val, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return fmt.Errorf("cannot scan %q into SmallInt: %w", v, err)
		}
		*si = SmallInt(val)
	case []byte:
		return si.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into SmallInt", val)
	}
	return nil
}

// Value implements driver.Valuer interface for SmallInt.
// driver.Value only allows int64 (not int32), per database/sql/driver contract.
func (si SmallInt) Value() (driver.Value, error) {
	return int64(si), nil
}

// NullSmallInt represents a nullable DB2 SMALLINT value.
type NullSmallInt struct {
	SmallInt SmallInt
	Valid    bool
}

// Scan implements sql.Scanner interface for NullSmallInt.
func (ns *NullSmallInt) Scan(val interface{}) error {
	if val == nil {
		ns.SmallInt = 0
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return ns.SmallInt.Scan(val)
}

// Value implements driver.Valuer interface for NullSmallInt.
func (ns NullSmallInt) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.SmallInt.Value()
}

// Int32 converts NullSmallInt to *int32 (nil if not valid).
func (ns NullSmallInt) Int32() *int32 {
	if !ns.Valid {
		return nil
	}
	v := int32(ns.SmallInt)
	return &v
}

// SQLNullInt32 is an alias for the standard sql.NullInt32
// provided for consistency with custom types above.
// Use this for generic integer fields that may be NULL.
type SQLNullInt32 = sql.NullInt32

// Date represents a DB2 DATE column (no time-of-day component).
// The driver binds plain time.Time values as a full TIMESTAMP structure,
// which DB2 rejects (SQL0180N) for DATE columns. Date implements
// driver.Valuer to send a "YYYY-MM-DD" string instead.
//
// Usage in struct:
//
//	type Event struct {
//	    EventDate db2dialect.Date `bun:"type:DATE"`
//	}
type Date time.Time

// Scan implements sql.Scanner interface for Date.
func (d *Date) Scan(val interface{}) error {
	if val == nil {
		*d = Date(time.Time{})
		return nil
	}
	switch v := val.(type) {
	case time.Time:
		*d = Date(v)
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("cannot scan %q into Date: %w", v, err)
		}
		*d = Date(t)
	case []byte:
		return d.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into Date", val)
	}
	return nil
}

// Value implements driver.Valuer interface for Date, formatting as "YYYY-MM-DD".
func (d Date) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

// Time converts Date to native time.Time.
func (d Date) Time() time.Time {
	return time.Time(d)
}

// TimeOfDay represents a DB2 TIME column (no date component).
// The driver binds plain time.Time values as a full TIMESTAMP structure,
// which DB2 rejects (SQL0180N) for TIME columns. TimeOfDay implements
// driver.Valuer to send a "HH:MM:SS" string instead.
//
// Usage in struct:
//
//	type Schedule struct {
//	    StartTime db2dialect.TimeOfDay `bun:"type:TIME"`
//	}
type TimeOfDay time.Time

// Scan implements sql.Scanner interface for TimeOfDay.
func (t *TimeOfDay) Scan(val interface{}) error {
	if val == nil {
		*t = TimeOfDay(time.Time{})
		return nil
	}
	switch v := val.(type) {
	case time.Time:
		*t = TimeOfDay(v)
	case string:
		parsed, err := time.Parse("15:04:05", v)
		if err != nil {
			return fmt.Errorf("cannot scan %q into TimeOfDay: %w", v, err)
		}
		*t = TimeOfDay(parsed)
	case []byte:
		return t.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into TimeOfDay", val)
	}
	return nil
}

// Value implements driver.Valuer interface for TimeOfDay, formatting as "HH:MM:SS".
func (t TimeOfDay) Value() (driver.Value, error) {
	return time.Time(t).Format("15:04:05"), nil
}

// Time converts TimeOfDay to native time.Time.
func (t TimeOfDay) Time() time.Time {
	return time.Time(t)
}

// Timestamp represents a DB2 TIMESTAMP column.
// The driver binds plain time.Time values via a SQL_TIMESTAMP_STRUCT C type,
// which can fail with SQL0180N depending on the column's fractional-seconds
// precision. Timestamp implements driver.Valuer to send a
// "YYYY-MM-DD HH:MM:SS.ffffff" string instead, which DB2 CLI parses reliably.
//
// Usage in struct:
//
//	type Project struct {
//	    StartDate db2dialect.Timestamp `bun:""`
//	}
type Timestamp time.Time

// Scan implements sql.Scanner interface for Timestamp.
func (ts *Timestamp) Scan(val interface{}) error {
	if val == nil {
		*ts = Timestamp(time.Time{})
		return nil
	}
	switch v := val.(type) {
	case time.Time:
		*ts = Timestamp(v)
	case string:
		t, err := time.Parse("2006-01-02 15:04:05.000000", v)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", v)
			if err != nil {
				return fmt.Errorf("cannot scan %q into Timestamp: %w", v, err)
			}
		}
		*ts = Timestamp(t)
	case []byte:
		return ts.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into Timestamp", val)
	}
	return nil
}

// Value implements driver.Valuer interface for Timestamp,
// formatting as "YYYY-MM-DD HH:MM:SS.ffffff".
func (ts Timestamp) Value() (driver.Value, error) {
	return time.Time(ts).Format("2006-01-02 15:04:05.000000"), nil
}

// Time converts Timestamp to native time.Time.
func (ts Timestamp) Time() time.Time {
	return time.Time(ts)
}
