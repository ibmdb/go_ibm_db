package main

import (
        "fmt"
        "strings"
        "testing"
)

// Issue 116
func TestDecimalColumn(t *testing.T) {
	if DecimalColumn() != 1 {
		t.Error("Error at DecimalColumn")
	}
}

func DecimalColumn() int {
        var tableOne string= "godecmaltable"

        db := Createconnection()
	defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Query("CREATE table " + tableOne + "(col1 DECIMAL(30, 2))")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Query error: ", err)
                return 0
        }

        _, err = db.Query("INSERT into " + tableOne + "(col1) values(9999999999999999999999999999.99)")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Query error: ", err)
                return 0
        }

        _, err = db.Query("INSERT into " + tableOne + "(col1) values(99999999999999999999)")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Query error: ", err)
                return 0
        }

        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                fmt.Println("Query error: ", err)
                return 0
        }

        defer rows.Close()
        for rows.Next() {
              var f  string
              err = rows.Scan(&f)
              if err != nil {
                      fmt.Println("Scan error: ", err)
                      return 0
              }
              //fmt.Printf("%v \n", f)
        }

        db.Query("DROP table " + tableOne)

        return 1
}

