// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"reflect"
	"time"

	"github.com/uptrace/bun/schema"
)

// OnTable processes and normalizes DB2 types for a table schema
// Maps Go types to appropriate DB2 SQL types
func (d *Dialect) OnTable(table *schema.Table) {
	for _, field := range table.Fields {
		// Skip already configured types
		if field.CreateTableSQLType != "" {
			continue
		}

		// Map Go types to DB2 SQL types
		d.mapFieldType(field)
	}
}

// mapFieldType determines the appropriate DB2 SQL type for a field
func (d *Dialect) mapFieldType(field *schema.Field) {
	// Get the underlying Go type
	typ := field.StructField.Type

	// Dereference pointer types
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Special handling for the custom SMALLINT helper types, which would
	// otherwise be misclassified as INTEGER (int32 kind) or VARCHAR (struct kind)
	switch typ {
	case reflect.TypeOf(SmallInt(0)), reflect.TypeOf(SmallIntBool(0)),
		reflect.TypeOf(NullSmallInt{}), reflect.TypeOf(NullSmallIntBool{}):
		field.CreateTableSQLType = "SMALLINT"
		return
	}

	// Map based on Go type
	switch typ.Kind() {
	case reflect.Bool:
		// DB2: Use SMALLINT for boolean (0 or 1)
		field.CreateTableSQLType = "SMALLINT"

	case reflect.String:
		// DB2: Use VARCHAR
		field.CreateTableSQLType = "VARCHAR(255)"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		// DB2: Use INTEGER for int types
		field.CreateTableSQLType = "INTEGER"

	case reflect.Int64:
		// DB2: Use BIGINT for int64
		field.CreateTableSQLType = "BIGINT"

	case reflect.Float32:
		// DB2: Use REAL for float32
		field.CreateTableSQLType = "REAL"

	case reflect.Float64:
		// DB2: Use DOUBLE PRECISION for float64
		field.CreateTableSQLType = "DOUBLE PRECISION"

	case reflect.Slice:
		// Check if it's a byte slice ([]byte)
		if typ.Elem().Kind() == reflect.Uint8 {
			field.CreateTableSQLType = "BLOB"
		}

	case reflect.Struct:
		// Special handling for time.Time and the custom time-based types
		switch typ {
		case reflect.TypeOf(time.Time{}):
			field.CreateTableSQLType = "TIMESTAMP"
		case reflect.TypeOf(Date{}):
			field.CreateTableSQLType = "DATE"
		case reflect.TypeOf(TimeOfDay{}):
			field.CreateTableSQLType = "TIME"
		case reflect.TypeOf(Timestamp{}):
			field.CreateTableSQLType = "TIMESTAMP"
		default:
			// Catch any other named type whose underlying type is time.Time
			if typ.ConvertibleTo(reflect.TypeOf(time.Time{})) {
				field.CreateTableSQLType = "TIMESTAMP"
			}
		}
	}
}
