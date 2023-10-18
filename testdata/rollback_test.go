package main

import (
        "testing"
)

func TestRollbackTransaction(t *testing.T) {
	if RollbackTransaction() != nil {
		t.Error("Error in Rollback Transaction")
	}
}


func RollbackTransaction() error {
        db := Createconnection()
        defer db.Close()

        bg, err := db.Begin()
        if err != nil {
                return err
        }

        _, err = bg.Exec("CREATE table gorollback(C1 int, C2 float, C3 double, C4 char, C5 varchar(30))")
        if err != nil {
                return err
        }

        err = bg.Rollback()
        if err != nil {
                return err
        }

        return nil
}

