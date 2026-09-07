// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/uptrace/bun/schema"
)

// TestMapFieldTypeBasicTypes verifies basic Go type mappings to DB2 types
func TestMapFieldTypeBasicTypes(t *testing.T) {
	d := New()

	tests := []struct {
		name     string
		goType   reflect.Type
		expected string
	}{
		{"bool", reflect.TypeOf(true), "SMALLINT"},
		{"string", reflect.TypeOf(""), "VARCHAR(255)"},
		{"int", reflect.TypeOf(int(0)), "INTEGER"},
		{"int32", reflect.TypeOf(int32(0)), "INTEGER"},
		{"int64", reflect.TypeOf(int64(0)), "BIGINT"},
		{"float32", reflect.TypeOf(float32(0)), "REAL"},
		{"float64", reflect.TypeOf(float64(0)), "DOUBLE PRECISION"},
		{"[]byte", reflect.TypeOf([]byte{}), "BLOB"},
		{"time.Time", reflect.TypeOf(time.Time{}), "TIMESTAMP"},
		{"Date", reflect.TypeOf(Date{}), "DATE"},
		{"TimeOfDay", reflect.TypeOf(TimeOfDay{}), "TIME"},
		{"Timestamp", reflect.TypeOf(Timestamp{}), "TIMESTAMP"},
		{"SmallInt", reflect.TypeOf(SmallInt(0)), "SMALLINT"},
		{"SmallIntBool", reflect.TypeOf(SmallIntBool(0)), "SMALLINT"},
		{"NullSmallInt", reflect.TypeOf(NullSmallInt{}), "SMALLINT"},
		{"NullSmallIntBool", reflect.TypeOf(NullSmallIntBool{}), "SMALLINT"},
		{"sql.NullBool", reflect.TypeOf(sql.NullBool{}), "SMALLINT"},
		{"sql.NullString", reflect.TypeOf(sql.NullString{}), "VARCHAR(255)"},
		{"sql.NullInt64", reflect.TypeOf(sql.NullInt64{}), "BIGINT"},
		{"sql.NullFloat64", reflect.TypeOf(sql.NullFloat64{}), "DOUBLE PRECISION"},
		{"sql.NullTime", reflect.TypeOf(sql.NullTime{}), "TIMESTAMP"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock field for testing
			field := &schema.Field{
				StructField: reflect.StructField{
					Type: test.goType,
				},
			}

			d.mapFieldType(field)

			if field.DiscoveredSQLType != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, field.DiscoveredSQLType)
			} else {
				t.Logf("✓ %s -> %s", test.name, field.DiscoveredSQLType)
			}
		})
	}
}

// TestMapFieldTypePointerTypes verifies pointer type dereferencing
func TestMapFieldTypePointerTypes(t *testing.T) {
	d := New()

	tests := []struct {
		name     string
		goType   reflect.Type
		expected string
	}{
		{"*string", reflect.TypeOf((*string)(nil)), "VARCHAR(255)"},
		{"*int64", reflect.TypeOf((*int64)(nil)), "BIGINT"},
		{"*bool", reflect.TypeOf((*bool)(nil)), "SMALLINT"},
		{"*time.Time", reflect.TypeOf((*time.Time)(nil)), "TIMESTAMP"},
		{"*Date", reflect.TypeOf((*Date)(nil)), "DATE"},
		{"*TimeOfDay", reflect.TypeOf((*TimeOfDay)(nil)), "TIME"},
		{"*Timestamp", reflect.TypeOf((*Timestamp)(nil)), "TIMESTAMP"},
		{"*SmallInt", reflect.TypeOf((*SmallInt)(nil)), "SMALLINT"},
		{"*SmallIntBool", reflect.TypeOf((*SmallIntBool)(nil)), "SMALLINT"},
		{"*NullSmallInt", reflect.TypeOf((*NullSmallInt)(nil)), "SMALLINT"},
		{"*NullSmallIntBool", reflect.TypeOf((*NullSmallIntBool)(nil)), "SMALLINT"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			field := &schema.Field{
				StructField: reflect.StructField{
					Type: test.goType,
				},
			}

			d.mapFieldType(field)

			if field.DiscoveredSQLType != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, field.DiscoveredSQLType)
			} else {
				t.Logf("✓ Pointer %s -> %s", test.name, field.DiscoveredSQLType)
			}
		})
	}
}

// TestOnTablePreservesUserSQLType verifies an explicit `type:` tag is not overwritten.
// Bun resolves CreateTableSQLType from UserSQLType/DiscoveredSQLType only after OnTable returns.
func TestOnTablePreservesUserSQLType(t *testing.T) {
	d := New()

	tagged := &schema.Field{
		StructField: reflect.StructField{Type: reflect.TypeOf("")},
		UserSQLType: "VARCHAR(100)",
	}
	untagged := &schema.Field{
		StructField: reflect.StructField{Type: reflect.TypeOf("")},
	}

	d.OnTable(&schema.Table{Fields: []*schema.Field{tagged, untagged}})

	if tagged.DiscoveredSQLType != "" {
		t.Errorf("user-tagged field should be left untouched, got DiscoveredSQLType=%q", tagged.DiscoveredSQLType)
	}
	if tagged.UserSQLType != "VARCHAR(100)" {
		t.Errorf("expected UserSQLType %q, got %q", "VARCHAR(100)", tagged.UserSQLType)
	}
	if untagged.DiscoveredSQLType != "VARCHAR(255)" {
		t.Errorf("expected untagged field to map to VARCHAR(255), got %q", untagged.DiscoveredSQLType)
	}
}

// TestValuersReturnDriverValueTypes verifies Value() returns types allowed by
// database/sql/driver (int64, not int32).
func TestValuersReturnDriverValueTypes(t *testing.T) {
	valuers := []struct {
		name string
		v    driver.Valuer
	}{
		{"SmallIntBool", SmallIntBool(1)},
		{"SmallInt", SmallInt(7)},
		{"NullSmallIntBool", NullSmallIntBool{SmallIntBool: 1, Valid: true}},
		{"NullSmallInt", NullSmallInt{SmallInt: 7, Valid: true}},
	}

	for _, test := range valuers {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.v.Value()
			if err != nil {
				t.Fatalf("Value() returned error: %v", err)
			}
			if !driver.IsValue(got) {
				t.Errorf("Value() returned non-driver.Value type %T", got)
			}
			if _, ok := got.(int64); !ok {
				t.Errorf("expected int64, got %T", got)
			}
		})
	}
}

// TestNullValuersReturnNilWhenInvalid verifies nullable helpers emit SQL NULL.
func TestNullValuersReturnNilWhenInvalid(t *testing.T) {
	for _, v := range []driver.Valuer{NullSmallIntBool{}, NullSmallInt{}} {
		got, err := v.Value()
		if err != nil {
			t.Fatalf("Value() returned error: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil for invalid %T, got %v", v, got)
		}
	}
}

func TestNullScansRemainInvalidAfterError(t *testing.T) {
	boolValue := NullSmallIntBool{SmallIntBool: 1}
	if err := boolValue.Scan("not-a-number"); err == nil {
		t.Fatal("expected invalid boolean scan to fail")
	}
	if boolValue.Valid {
		t.Fatal("boolean value should remain invalid after a failed scan")
	}
	if got, err := boolValue.Value(); err != nil || got != nil {
		t.Fatalf("failed boolean scan should emit NULL, got value=%v err=%v", got, err)
	}

	intValue := NullSmallInt{SmallInt: 7}
	if err := intValue.Scan("not-a-number"); err == nil {
		t.Fatal("expected invalid integer scan to fail")
	}
	if intValue.Valid {
		t.Fatal("integer value should remain invalid after a failed scan")
	}
	if got, err := intValue.Value(); err != nil || got != nil {
		t.Fatalf("failed integer scan should emit NULL, got value=%v err=%v", got, err)
	}
}
