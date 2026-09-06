// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"database/sql"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/feature"
	"github.com/uptrace/bun/schema"
)

// DB2 dialect name constant
const db2Name dialect.Name = 10 // Custom value for DB2 dialect

// Dialect represents the DB2 SQL dialect for Bun ORM
type Dialect struct {
	schema.BaseDialect
	tables   *schema.Tables  // Registry for table schemas
	features feature.Feature // Enabled DB2 features
}

// New creates a new DB2 dialect instance
func New() *Dialect {
	d := &Dialect{
		// Identity is required so Bun emits the GENERATED ... AS IDENTITY clause
		// (via AppendSequence) and omits autoincrement PK columns from INSERT statements.
		features: feature.CTE | feature.WithValues | feature.SelectExists | feature.CompositeIn | feature.OffsetFetch | feature.Identity,
	}
	d.tables = schema.NewTables(d)
	return d
}

// IdentQuote returns the identifier quote character (double quote for DB2)
func (d *Dialect) IdentQuote() byte {
	return '"'
}

// Name returns the dialect name for DB2
func (d *Dialect) Name() dialect.Name {
	return db2Name
}

// Features returns the enabled DB2 dialect features
func (d *Dialect) Features() feature.Feature {
	return d.features
}

// Tables returns the table schema registry
func (d *Dialect) Tables() *schema.Tables {
	if d.tables == nil {
		d.tables = schema.NewTables(d)
	}
	return d.tables
}

// DefaultVarcharLen returns the default VARCHAR length for DB2
func (d *Dialect) DefaultVarcharLen() int {
	return 255
}

// DefaultSchema returns an empty schema so DB2 uses the current authorization ID.
func (d *Dialect) DefaultSchema() string {
	return ""
}

// Init initializes the dialect with a database connection
// This is called by Bun during initialization
func (d *Dialect) Init(db *sql.DB) {
	// Placeholder for future initialization logic
	// Currently, lazy initialization in Tables() is sufficient
}
