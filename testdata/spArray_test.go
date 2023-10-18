package main

import (
        //"database/sql"
        "fmt"
        "strings"
        "testing"
)

func TestStoredProcedureArray(t *testing.T) {
	if StoredProcedureArray_1() != nil {
		t.Error("Error at stored procedure array")
	}
}

func StoredProcedureArray_1() error{
        db := Createconnection()

        var arr1 = []int{10, 20, 30, 40, 50}
        var arr2 = []string{"Row 10", "Row 20", "Row 30", "Row 40", "Row 50"}

        db.Exec("Drop table TT")

        _, err1 := db.Exec("Drop procedure sp1(INTEGER, VARCHAR(10))")
        if err1 !=  nil {
                fmt.Println("Drop procedure error : ", err1)
                return err1
        }

        _, err := db.Exec("create table TT(C1 INTEGER NOT NULL, C2 VARCHAR(10))");
        if err != nil {
                fmt.Println("Error: ", err)
                return err
        }

        st, err := db.Prepare("create procedure sp1(in arr1 INTEGER, in arr2 VARCHAR(10)) LANGUAGE SQL BEGIN  INSERT INTO TT VALUES(arr1, arr2); END")
        if err != nil {
                fmt.Println("Error: ", err)
                return err
        }

        _, err = st.Query()
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                return err
        }

        _, err = db.Exec("CALL sp1(?,?)", arr1, arr2)
        if err != nil {
                fmt.Println("Error: ", err)
                return err
        }
/*
        queryTypes := "select * from TT;"

        rows, err := db.Query(queryTypes)
        if err != nil {
                fmt.Println("Error:", err)
                return err
        }

        for rows.Next() {
                var c1 int
                var c2 string
                err = rows.Scan(&c1, &c2)
                if err != nil {
                        fmt.Println("Error:", err)
                }
                fmt.Printf("C1=: %v\t C2: %v\n", c1, c2)
        }
*/
        return nil
}

