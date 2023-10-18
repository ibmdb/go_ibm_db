package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestInt2Array(t *testing.T) {
	if Int2Array() != nil {
		t.Error("Error at Int2Array")
	}
}

// A large integer is binary integer with a precision of 31 bits. The range is -2147483648 to +2147483647.
func Int2Array() error {
        var tableOne string= "goarr"
        var errStr string

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 int)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []int{-2147483648, -32769, 0, 32768,  2147483647}
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []int2")
                return err
        }

        substring := "SQLSTATE=22003"
        c :=  []int{6}
        d :=  []int{-2147483649}
         st, err = db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(c, d)
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
                if !strings.Contains(errStr, substring) {
                        return err
                }
        }

        e :=  []int{7}
        f :=  []int{2147483648}
         st, err = db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(e, f)
        if err != nil {
                errStr = fmt.Sprintf("%s", err)
                if !strings.Contains(errStr, substring) {
                        return err
                }
        }

        rows, err2 := db.Query("SELECT * from " + tableOne)
        if err2 != nil {
                return err
        }

        defer rows.Close()
        for rows.Next() {
              var c1, c2  string
              err = rows.Scan(&c1, &c2)
              if err != nil {
                      return err
              }

              //fmt.Printf("%v  %v \n", c1, c2)
        }

        db.Query("DROP table " + tableOne)

        return nil
}

