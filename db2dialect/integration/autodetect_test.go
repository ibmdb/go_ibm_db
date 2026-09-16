package main

import (
	"context"
	"strings"
	"testing"

	go_ibm_db "github.com/ibmdb/go_ibm_db"
	"github.com/ibmdb/go_ibm_db/api"
	"github.com/ibmdb/go_ibm_db/db2dialect"
)

func TestBun_AutoDetectPlatform(t *testing.T) {
	sqlDB := OpenTestDB(t)
	defer sqlDB.Close()

	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	var dbmsName string
	if err := conn.Raw(func(driverConn any) error {
		var rawErr error
		dbmsName, rawErr = driverConn.(*go_ibm_db.Conn).GetInfo(api.SQL_DBMS_NAME)
		return rawErr
	}); err != nil {
		t.Fatal(err)
	}
	t.Logf("SQL_DBMS_NAME: %s", dbmsName)

	wantTarget := targetForDBMSName(dbmsName)
	dialect := db2dialect.New()
	dialect.Init(sqlDB)

	if dialect.Target() != wantTarget {
		t.Fatalf("auto-detected target %v, want %v (SQL_DBMS_NAME=%s)", dialect.Target(), wantTarget, dbmsName)
	}
}

func targetForDBMSName(name string) db2dialect.TargetPlatform {
	upper := strings.ToUpper(name)
	switch {
	case upper == "DB2" || strings.HasPrefix(upper, "DSN"):
		return db2dialect.TargetZOS
	case strings.HasPrefix(upper, "AS"):
		return db2dialect.TargetIBMi
	default:
		return db2dialect.TargetLUW
	}
}
