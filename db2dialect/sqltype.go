// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/uptrace/bun/schema"
)

// OnTable processes and normalizes DB2 types for a table schema
// Maps Go types to appropriate DB2 SQL types
func (d *Dialect) OnTable(table *schema.Table) {
	for _, field := range table.Fields {
		// On z/OS, 'DEFAULT current_timestamp' in CREATE TABLE produces SQL0104N/SQL0199N/SQL0637N syntax error.
		// Clear SQLDefault on z/OS so Bun emits 'TIMESTAMP' cleanly without an invalid DEFAULT clause.
		if d.target == TargetZOS && (field.SQLDefault == "current_timestamp" || field.SQLDefault == "CURRENT_TIMESTAMP" || field.SQLDefault == "WITH DEFAULT") {
			field.SQLDefault = ""
		}

		// Bun resolves CreateTableSQLType from UserSQLType/DiscoveredSQLType only after
		// OnTable returns, so an explicit `type:` tag must be preserved here.
		if field.UserSQLType != "" || field.CreateTableSQLType != "" {
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
		field.DiscoveredSQLType = "SMALLINT"
		return
	}

	// Map based on Go type
	switch typ.Kind() {
	case reflect.Bool:
		// DB2: Use SMALLINT for boolean (0 or 1)
		field.DiscoveredSQLType = "SMALLINT"

	case reflect.String:
		// DB2: Use VARCHAR
		field.DiscoveredSQLType = "VARCHAR(255)"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		// DB2: Use INTEGER for int types
		field.DiscoveredSQLType = "INTEGER"

	case reflect.Int64:
		// DB2: Use BIGINT for int64
		field.DiscoveredSQLType = "BIGINT"

	case reflect.Float32:
		// DB2: Use REAL for float32
		field.DiscoveredSQLType = "REAL"

	case reflect.Float64:
		// DB2: Use DOUBLE PRECISION for float64
		field.DiscoveredSQLType = "DOUBLE PRECISION"

	case reflect.Slice:
		// Check if it's a byte slice ([]byte)
		if typ.Elem().Kind() == reflect.Uint8 {
			field.DiscoveredSQLType = "BLOB"
		}

	case reflect.Struct:
		// Special handling for time.Time, database/sql Null types, and custom time-based types
		switch typ {
		case reflect.TypeOf(time.Time{}):
			field.DiscoveredSQLType = "TIMESTAMP"
		case reflect.TypeOf(Date{}):
			field.DiscoveredSQLType = "DATE"
		case reflect.TypeOf(TimeOfDay{}):
			field.DiscoveredSQLType = "TIME"
		case reflect.TypeOf(Timestamp{}):
			field.DiscoveredSQLType = "TIMESTAMP"
		case reflect.TypeOf(sql.NullBool{}):
			field.DiscoveredSQLType = "SMALLINT"
		case reflect.TypeOf(sql.NullString{}):
			field.DiscoveredSQLType = "VARCHAR(255)"
		case reflect.TypeOf(sql.NullInt64{}):
			field.DiscoveredSQLType = "BIGINT"
		case reflect.TypeOf(sql.NullInt32{}):
			field.DiscoveredSQLType = "INTEGER"
		case reflect.TypeOf(sql.NullInt16{}):
			field.DiscoveredSQLType = "SMALLINT"
		case reflect.TypeOf(sql.NullByte{}):
			field.DiscoveredSQLType = "SMALLINT"
		case reflect.TypeOf(sql.NullFloat64{}):
			field.DiscoveredSQLType = "DOUBLE PRECISION"
		case reflect.TypeOf(sql.NullTime{}):
			field.DiscoveredSQLType = "TIMESTAMP"
		default:
			// Catch any other named type whose underlying type is time.Time
			if typ.ConvertibleTo(reflect.TypeOf(time.Time{})) {
				field.DiscoveredSQLType = "TIMESTAMP"
			}
		}
	}
}
