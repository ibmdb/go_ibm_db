package main

import (
	"database/sql"
        "fmt"
        "testing"
)

func TestStoredProcedureInOut(t *testing.T) {
	if StoredProcedureInOut() != nil {
		t.Error("Error at stored procedure")
	}
}

//StoredProcedureInOut function tests OUT Parameter by calling get_dbsize_info.
func StoredProcedureInOut() error {
        db := Createconnection()
        defer db.Close()
        in1 := 10
        inout1 := 2
        var out1, out2 int
        st, err := db.Prepare("create or replace procedure sp2(in var1 integer, inout var2 integer, out var3 integer, out var4 integer) LANGUAGE SQL BEGIN   SET var2 = var1 + var2; SET var3 = var1 - var2; SET var4 = var1 * var2; END")
        if err != nil {
                return err
        }
        st.Query()
        _, err = db.Exec("CALL sp2(?,?,?,?)", in1, sql.Out{Dest: &inout1, In: true}, sql.Out{Dest: &out1}, sql.Out{Dest: &out2})
        if err != nil {
                return err
        }
        if inout1 != 12 || out1 != -2 || out2 != 120 {
                return fmt.Errorf("Wrong data retrieved")
        }
        return nil
}


