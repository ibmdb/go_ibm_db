// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/feature"
	"github.com/uptrace/bun/schema"
)

// errUnsupportedDriverConn is returned when the underlying *sql.DB's driver
// connection does not implement dbmsNamer (e.g. it isn't go_ibm_db).
var errUnsupportedDriverConn = errors.New("db2dialect: driver connection does not support DBMSName() detection")

// dbmsNamer is satisfied by go_ibm_db's *Conn via structural typing, so
// db2dialect never needs to import github.com/ibmdb/go_ibm_db directly.
type dbmsNamer interface {
	DBMSName() (string, error)
}

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

// classifyDBMSName maps an ODBC SQL_DBMS_NAME string to a TargetPlatform,
// matching the heuristic used by IBM's own python-ibmdb test suite:
// z/OS reports the exact string "DB2" (no platform suffix) or a name
// prefixed with "DSN"; IBM i is prefixed "AS"; Informix is prefixed "IDS"
// (unsupported by this driver, treated as unrecognized); LUW is prefixed
// "DB2/" (e.g. "DB2/LINUXX8664"). Anything else defaults to LUW.
func classifyDBMSName(name string) TargetPlatform {
	upper := strings.ToUpper(name)
	switch {
	case upper == "DB2" || strings.HasPrefix(upper, "DSN"):
		return TargetZOS
	case strings.HasPrefix(upper, "AS"):
		return TargetIBMi
	case strings.HasPrefix(upper, "DB2/"):
		return TargetLUW
	default:
		log.Printf("db2dialect: WARNING: unable to deduce DB2 platform from DBMS_NAME %q, defaulting to LUW", name)
		return TargetLUW
	}
}

// Init initializes the dialect with a database connection.
// Automatically detects the target platform flavor (LUW, z/OS, or IBM i) if
// not explicitly set, via a single SQLGetInfo(SQL_DBMS_NAME) call. If
// detection fails for any reason (unsupported driver connection, SQLGetInfo
// error, or an unrecognized DBMS_NAME), it logs a warning and defaults to
// LUW rather than failing; use NewLUW()/NewZOS()/NewIBMi() to skip detection.
func (d *Dialect) Init(db *sql.DB) {
	if db == nil || d.targetSetExplicitly || d.autoDetected {
		return
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Printf("db2dialect: WARNING: unable to deduce DB2 platform (%v), defaulting to LUW", err)
		d.target = TargetLUW
		d.autoDetected = true
		return
	}
	defer conn.Close()

	var name string
	rawErr := conn.Raw(func(driverConn any) error {
		namer, ok := driverConn.(dbmsNamer)
		if !ok {
			return errUnsupportedDriverConn
		}
		var callErr error
		name, callErr = namer.DBMSName()
		return callErr
	})
	if rawErr != nil {
		log.Printf("db2dialect: WARNING: unable to deduce DB2 platform (%v), defaulting to LUW", rawErr)
		d.target = TargetLUW
		d.autoDetected = true
		return
	}

	d.target = classifyDBMSName(name)
	d.autoDetected = true
}
