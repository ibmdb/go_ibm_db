package main

import (
	"os"
	"strings"
	"testing"

	"github.com/ibmdb/go_ibm_db/db2dialect"
	"github.com/uptrace/bun"
)

func TestBun_AutoDetectPlatform(t *testing.T) {
	sqlDB := OpenTestDB(t)
	defer sqlDB.Close()

	dialect := db2dialect.New()
	db := bun.NewDB(sqlDB, dialect)
	defer db.Close()

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
