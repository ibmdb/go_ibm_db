// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/feature"
	"github.com/uptrace/bun/schema"
)

// db2Name is intentionally outside Bun v1.2.18's built-in dialect range
// (Invalid through Oracle, values 0-5). Keep this value stable and re-check
// it when upgrading Bun so a newly added built-in dialect cannot collide.
const db2Name dialect.Name = 10

// TargetPlatform represents the target DB2 platform flavor
type TargetPlatform int

const (
	// TargetLUW targets DB2 for Linux, Unix, and Windows (default)
	TargetLUW TargetPlatform = iota
	// TargetZOS targets DB2 for z/OS (Mainframe)
	TargetZOS
	// TargetIBMi targets DB2 for IBM i (AS400/iSeries)
	TargetIBMi
)

// String returns the string representation of the target platform
func (t TargetPlatform) String() string {
	switch t {
	case TargetZOS:
		return "z/OS"
	case TargetIBMi:
		return "IBM i"
	default:
		return "LUW"
	}
}

// Option configures a Dialect instance
type Option func(*Dialect)

// WithTarget sets the target DB2 platform flavor
func WithTarget(target TargetPlatform) Option {
	return func(d *Dialect) {
		d.target = target
		d.targetSetExplicitly = true
	}
}

// Dialect represents the DB2 SQL dialect for Bun ORM
type Dialect struct {
	schema.BaseDialect
	tables              *schema.Tables  // Registry for table schemas
	features            feature.Feature // Enabled DB2 features
	target              TargetPlatform  // Target DB2 platform flavor
	targetSetExplicitly bool            // Whether target was set explicitly via WithTarget
	autoDetected        bool            // Whether target was auto-detected from connection
}

// New creates a new DB2 dialect instance with optional functional options
func New(opts ...Option) *Dialect {
	d := &Dialect{
		// Identity is required so Bun emits the GENERATED ... AS IDENTITY clause
		// (via AppendSequence) and omits autoincrement PK columns from INSERT statements.
		features: feature.CTE | feature.WithValues | feature.SelectExists | feature.CompositeIn | feature.OffsetFetch | feature.Identity,
		target:   TargetLUW, // LUW is default for backward compatibility
	}

	for _, opt := range opts {
		opt(d)
	}

	d.tables = schema.NewTables(d)
	return d
}

// NewLUW creates a dialect explicitly configured for DB2 for LUW.
func NewLUW() *Dialect {
	return New(WithTarget(TargetLUW))
}

// NewZOS creates a dialect explicitly configured for DB2 for z/OS.
func NewZOS() *Dialect {
	return New(WithTarget(TargetZOS))
}

// NewIBMi creates a dialect explicitly configured for DB2 for IBM i.
func NewIBMi() *Dialect {
	return New(WithTarget(TargetIBMi))
}

// Target returns the target DB2 platform flavor
func (d *Dialect) Target() TargetPlatform {
	return d.target
}

// DummyTable returns the dummy table name for non-table/dual queries.
// On z/OS, every SELECT statement strictly requires a FROM clause (SYSIBM.SYSDUMMY1).
func (d *Dialect) DummyTable() string {
	return "SYSIBM.SYSDUMMY1"
}

// CatalogSchema returns the system catalog schema name for the target platform.
// LUW uses SYSCAT, z/OS uses SYSIBM, and IBM i uses QSYS2.
func (d *Dialect) CatalogSchema() string {
	switch d.target {
	case TargetZOS:
		return "SYSIBM"
	case TargetIBMi:
		return "QSYS2"
	default:
		return "SYSCAT"
	}
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

// AppendTime formats time.Time into DB2 SQL timestamp format ('YYYY-MM-DD HH:MM:SS.ffffff')
// DB2 TIMESTAMP columns do not accept RFC3339 timezone offsets (e.g. '+00:00').
func (d *Dialect) AppendTime(b []byte, tm time.Time) []byte {
	if tm.IsZero() {
		return append(b, "NULL"...)
	}
	b = append(b, '\'')
	b = tm.UTC().AppendFormat(b, "2006-01-02 15:04:05.000000")
	b = append(b, '\'')
	return b
}

// Init initializes the dialect with a database connection.
// Automatically detects the target platform flavor (LUW, z/OS, or IBM i) if not explicitly set.
func (d *Dialect) Init(db *sql.DB) {
	if db == nil || d.targetSetExplicitly || d.autoDetected {
		return
	}

	var dummy int
	// Try LUW check (SYSCAT.TABLES exists only on DB2 LUW)
	if err := db.QueryRow("SELECT 1 FROM SYSCAT.TABLES FETCH FIRST 1 ROWS ONLY").Scan(&dummy); err == nil {
		d.target = TargetLUW
		d.autoDetected = true
		return
	}

	// Try z/OS check (SYSIBM.SYSTABLES exists on z/OS)
	if err := db.QueryRow("SELECT 1 FROM SYSIBM.SYSTABLES FETCH FIRST 1 ROWS ONLY").Scan(&dummy); err == nil {
		d.target = TargetZOS
		d.autoDetected = true
		return
	}

	// Try IBM i check (QSYS2.SYSTABLES exists on IBM i / AS400)
	if err := db.QueryRow("SELECT 1 FROM QSYS2.SYSTABLES FETCH FIRST 1 ROWS ONLY").Scan(&dummy); err == nil {
		d.target = TargetIBMi
		d.autoDetected = true
		return
	}
}
