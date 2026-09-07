package main

import (
	"os"
	"strings"

	"github.com/ibmdb/go_ibm_db/db2dialect"
)

// GetDialect returns the appropriate DB2 dialect for Bun ORM
// Configures target platform based on DB2_TARGET_PLATFORM env var (LUW, ZOS, IBMi)
func GetDialect() *db2dialect.Dialect {
	return getDialectImpl()
}

// getDialectImpl handles the actual dialect selection and target configuration
func getDialectImpl() *db2dialect.Dialect {
	targetEnv := strings.ToUpper(strings.TrimSpace(os.Getenv("DB2_TARGET_PLATFORM")))

	switch targetEnv {
	case "ZOS", "Z/OS", "MAINFRAME":
		return db2dialect.New(db2dialect.WithTarget(db2dialect.TargetZOS))
	case "IBMI", "AS400", "ISERIES":
		return db2dialect.New(db2dialect.WithTarget(db2dialect.TargetIBMi))
	case "LUW":
		return db2dialect.New(db2dialect.WithTarget(db2dialect.TargetLUW))
	default:
		// When no env var is set, return default dialect without explicit target,
		// allowing d.Init(db) to auto-detect the target platform from the active connection.
		return db2dialect.New()
	}
}
