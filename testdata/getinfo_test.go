package main

import (
	"strings"
	"testing"

	"github.com/ibmdb/go_ibm_db/api"
)

// TestGetInfo verifies SQLGetInfo support via Conn.GetInfo() for a few
// commonly used info types.
func TestGetInfo(t *testing.T) {
	infoTypes := map[string]api.SQLUSMALLINT{
		"SQL_DBMS_NAME":     api.SQL_DBMS_NAME,
		"SQL_DBMS_VER":      api.SQL_DBMS_VER,
		"SQL_DRIVER_NAME":   api.SQL_DRIVER_NAME,
		"SQL_DRIVER_VER":    api.SQL_DRIVER_VER,
		"SQL_DATABASE_NAME": api.SQL_DATABASE_NAME,
	}

	for label, infoType := range infoTypes {
		value, err := GetInfo(infoType)
		if err != nil {
			t.Errorf("GetInfo(%s) error: %v", label, err)
			continue
		}
		if strings.TrimSpace(value) == "" {
			t.Errorf("GetInfo(%s) returned an empty string", label)
			continue
		}
		t.Logf("%s: %s", label, value)
	}
}
