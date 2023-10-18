package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestVargraphicArray(t *testing.T) {
	if VargraphicArray_1() != nil {
		t.Error("Error at VargraphicArray")
	}
}

//VARGRAPHIC(n) Varying-length graphic strings. The maximum length, n, must be greater than 0
//and less than a number that depends on the page size of the table space. The maximum length is 16352.
func VargraphicArray_1() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 vargraphic(20))")
        if err != nil {
                return err
        }

        a :=  []int{1, 2, 3, 4, 5}
        b :=  []string  { "Hello World!!", "!@#$%^&*()", "1234567890", "How are you?", "Happy Birthday!" }
         st, err := db.Prepare("Insert into " +tableOne+ " values(?, ?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(a, b)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []vargraphic")
                return err
        }
        var errStr string
        substring := "SQLSTATE=22001"
        c :=  []int{6}
        //d :=  []string{"abcdefghijklmnopqurstuvwxyz"}
        d :=  []string{"123456789012345678901"}
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


