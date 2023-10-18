package main

import (
	"database/sql"
	"time"
	"testing"
)

func TestStoredProcedure(t *testing.T) {
	if StoredProcedure() != nil {
		t.Error("Error at stored procedure")
	}
}


//StoredProcedure function tests OUT Parameter by calling get_dbsize_info.
func StoredProcedure() error {
        var (
                snapTime   time.Time
                dbsize     int64
                dbcapacity int64
        )
        db := Createconnection()
        defer db.Close()
        _, err := db.Exec("call sysproc.get_dbsize_info(?, ?, ?,0)", sql.Out{Dest: &snapTime}, sql.Out{Dest: &dbsize}, sql.Out{Dest: &dbcapacity})
        if err != nil {
                return err
        }
        return nil
}

