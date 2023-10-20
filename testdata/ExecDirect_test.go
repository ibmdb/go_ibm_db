package main

import "testing"

func TestExecDirect(t *testing.T) {
	if ExecDirect() != nil {
		t.Error("Error in ExecDirect")
	}
}

//ExecDirect will execute the query without prepare
func ExecDirect() error {
        db := Createconnection()
        defer db.Close()
        _, err := db.Query("select * from rocket")
        if err != nil {
                return err
        }
        return nil
}

