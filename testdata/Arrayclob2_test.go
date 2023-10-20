package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestClob2Array(t *testing.T) {
	if Clob2Array_1() != nil {
		t.Error("Error at Clob2Array")
	}
}


//A CLOB (character large object) value can be up to 2,147,483,647 characters long.
//A CLOB is used to store unicode character-based data, such as large documents in any character set.
//The length is given in number characters for both CLOB, unless one of the suffixes K, M, or G is given,
//relating to the multiples of 1024, 1024*1024, 1024*1024*1024 respectively.
//{CLOB |CHARACTER LARGE OBJECT} [ ( length [{K |M |G}] ) ]
func Clob2Array_1() error {
        var tableOne string= "goarr"

        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 clob(64 K))")
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
                fmt.Println("Error while inserting []clob2")
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

