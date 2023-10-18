package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestBigintArray(t *testing.T) {
	if BigintArray() != nil {
		t.Error("Error at BigintArray")
	}
}


// A big integer is a binary integer with a precision of 63 bits. The range of big integers is -9223372036854775808 to +9223372036854775807.
func BigintArray() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 bigint)")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5, 6, 7}
        b :=  []int{-9223372036854775808, -2147483648, -32769, 0, 32768,  2147483647, 9223372036854775807}

        st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []bigint")
                return err
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

