// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
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

			if field.CreateTableSQLType != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, field.CreateTableSQLType)
			} else {
				t.Logf("✓ %s -> %s", test.name, field.CreateTableSQLType)
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

			if field.CreateTableSQLType != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, field.CreateTableSQLType)
			} else {
				t.Logf("✓ Pointer %s -> %s", test.name, field.CreateTableSQLType)
			}
		})
	}
}
