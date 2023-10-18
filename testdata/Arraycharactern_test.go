package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestCharacterArray(t *testing.T) {
	if CharacterArray() != nil {
		t.Error("Error at CharacterArray")
	}
}


//CHARACTER(n)Fixed-length character strings with a length of n bytes. n must be greater than 0 and not greater than 255. The default length is 1.
func CharacterArray() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 character(5))")
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
                fmt.Println("Error while inserting []character")
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

