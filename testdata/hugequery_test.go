package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestHugeQuery(t *testing.T) {
	if HugeQuery() != 1 {
		t.Error("Error at HugeQuery")
	}
}

func HugeQuery() int {
        var insertCount string  = "1"
        var tableOne string= "goleaktable1"
        var tableTwo string= "goleaktable2"

        var maxVarChar string = "LEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTESTLEAKTES"

        db := Createconnection()
	defer db.Close()

        db.Query("DROP table " + tableOne)
        db.Query("DROP table " + tableTwo)

        _, err := db.Query("CREATE table " + tableOne + "(PID VARCHAR(10), C1 VARCHAR(255), C2 VARCHAR(255), C3 VARCHAR(255))")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Query error: ", err)
                return 0
        }

        _, err = db.Query("CREATE table " + tableTwo + "(PID VARCHAR(10), C1 VARCHAR(255), C2 VARCHAR(255), C3 VARCHAR(255))")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Query error: ", err)
                return 0
        }
        query := "values('" + insertCount + "', '" + maxVarChar + "', '" + maxVarChar + "', '" + maxVarChar + "')"

        for i := 1; i <= 5; i++ {
                 _, err = db.Query("INSERT into " + tableOne + "(PID, C1, C2, C3) " + query)
                if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        fmt.Println("Query error: ", err)
                        return 0
                }

                _, err = db.Query("INSERT into " + tableTwo + "(PID, C1, C2, C3) " + query)
                if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                        fmt.Println("Query error: ", err)
                        return 0
                }
        }

        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                fmt.Println("Query error: ", err2)
                return 0
        }

        defer rows.Close()
        for rows.Next() {
              var t, x, m, n string
              err = rows.Scan(&t, &x, &m, &n)
              if err != nil {
                      fmt.Println("Scan error: ", err)
                      return 0
              }
              //fmt.Printf("%v  %v   %v  %v\n", t, x, m, n)
        }

        db.Query("DROP table " + tableOne)
        db.Query("DROP table " + tableTwo)

        return 1
}

