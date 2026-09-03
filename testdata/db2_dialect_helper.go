package main

import (
	"github.com/ibmdb/go_ibm_db/db2dialect"
)

// GetDialect returns the appropriate DB2 dialect for Bun ORM
//
// Dialect Selection Logic:
// 1. First tries to use uptrace's db2dialect (if available in future versions)
// 2. Falls back to local db2dialect implementation (currently used)
//
// To switch dialects, modify the getDialectImpl() function below
// No need to change imports throughout the codebase
func GetDialect() *db2dialect.Dialect {
	return getDialectImpl()
}

// getDialectImpl handles the actual dialect selection
// Modify this function to switch between uptrace and local dialects
func getDialectImpl() *db2dialect.Dialect {
	// Currently using local db2dialect implementation (from github.com/ibmdb/go_ibm_db/db2dialect)
	// This is the primary and recommended approach for these tests
	return db2dialect.New()

	// If uptrace ever provides db2dialect in future versions, you can switch to:
	// import "github.com/uptrace/bun/dialect/db2dialect" as bundb2
	// return bundb2.New()
}
