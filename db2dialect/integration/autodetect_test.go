package main

import (
	"os"
	"strings"
	"testing"

	"github.com/ibmdb/go_ibm_db/db2dialect"
	"github.com/uptrace/bun"
)

// TestBun_AutoDetectPlatform verifies db2dialect.New() (no explicit target)
// auto-detects the server platform via SQLGetInfo(SQL_DBMS_NAME) and matches
// DB2_TARGET_PLATFORM when that env var is set.
func TestBun_AutoDetectPlatform(t *testing.T) {
	sqlDB := OpenTestDB(t)
	defer sqlDB.Close()

	db := bun.NewDB(sqlDB, db2dialect.New())
	defer db.Close()

	dialect, ok := db.Dialect().(*db2dialect.Dialect)
	if !ok {
		t.Fatalf("expected *db2dialect.Dialect, got %T", db.Dialect())
	}

	want := strings.ToUpper(strings.TrimSpace(os.Getenv("DB2_TARGET_PLATFORM")))
	if want == "" {
		t.Logf("DB2_TARGET_PLATFORM not set; detected target: %v", dialect.Target())
		return
	}

	var wantTarget db2dialect.TargetPlatform
	switch want {
	case "ZOS", "Z/OS", "MAINFRAME":
		wantTarget = db2dialect.TargetZOS
	case "IBMI", "AS400", "ISERIES":
		wantTarget = db2dialect.TargetIBMi
	default:
		wantTarget = db2dialect.TargetLUW
	}

	if dialect.Target() != wantTarget {
		t.Fatalf("auto-detected target %v, want %v (DB2_TARGET_PLATFORM=%s)", dialect.Target(), wantTarget, want)
	}
}
