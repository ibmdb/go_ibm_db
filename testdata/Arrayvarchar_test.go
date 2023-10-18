package main

import (
        "fmt"
        "strings"
        "testing"
)


func TestVarcharArray(t *testing.T) {
	if VarcharArray_1() != nil {
		t.Error("Error at VarcharArray")
	}
}


//VARCHAR(n):Varying-length character strings with a maximum length of n bytes.
//n must be greater than 0 and less than a number that depends on the page size of the table space. The maximum length is 32704.
func VarcharArray_1() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 varchar(5))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "a", "ab", "abc", "abcd", "abcde" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []varchar")
                return err
        }

        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        d :=  []string{"abcdef"}
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

