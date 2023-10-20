package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestDecfloat(t *testing.T) {
	if Decfloat() != nil {
		t.Error("Error at Decfloat")
	}
}


func Decfloat() error {
        var tableOne string= "goarr"
        db := Createconnection()
        defer db.Close()

        db.Query("DROP table " + tableOne)

        _, err := db.Exec("CREATE table " + tableOne + "(col1 int, col2 decfloat)")
        if err != nil {
                return err
        }

        st1, err1 := db.Prepare("INSERT into " + tableOne + " values(1, 45.678)")
        defer st1.Close()
        if err1 != nil {
                return err1
        }
        _, err1 = st1.Query()
        if !strings.Contains(fmt.Sprint(err1), "did not create a result set") {
                return err1
        }

        st2, err2 := db.Prepare("INSERT into " + tableOne + " values(2, 0.2345600)")
        defer st2.Close()
        if err2 != nil {
                return err2
        }
        _, err2 = st2.Query()
        if !strings.Contains(fmt.Sprint(err2), "did not create a result set") {
                return err2
        }

        st3, err3 := db.Prepare("INSERT into " + tableOne + " values(3, 111e99)")
        defer st3.Close()
        if err3 != nil {
                return err3
        }
        _, err3 = st3.Query()
        if !strings.Contains(fmt.Sprint(err3), "did not create a result set") {
                return err3
        }


        st4, err4 := db.Prepare("INSERT into " + tableOne + " values(4, 111e-99)")
        defer st4.Close()
        if err4 != nil {
                return err4
        }
        _, err4 = st4.Query()
        if !strings.Contains(fmt.Sprint(err4), "did not create a result set") {
                return err4
        }

        st5, err5 := db.Prepare("INSERT into " + tableOne + " values(5, 100.2001234)")
        defer st5.Close()
        if err5 != nil {
                return err5
        }
       _, err5 = st5.Query()
        if !strings.Contains(fmt.Sprint(err5), "did not create a result set") {
                return err5
        }

        st6, err6 := db.Prepare("INSERT into " + tableOne + " values(6, -1000)")
        defer st6.Close()
        if err6 != nil {
                return err6
        }
        _, err6 = st6.Query()
        if !strings.Contains(fmt.Sprint(err6), "did not create a result set") {
                return err6
        }

        st7, err7 := db.Prepare("INSERT into " + tableOne + " values(7, -Inf)")
        defer st7.Close()
        if err7 != nil {
                return err7
        }
        _, err7 = st7.Query()
        if !strings.Contains(fmt.Sprint(err7), "did not create a result set") {
                return err7
        }

        rows, err8 := db.Query("SELECT * from " + tableOne)
        if err8 != nil {
                return err8
        }

        defer rows.Close()
        for rows.Next() {
              var c1, c2  string
              err9 := rows.Scan(&c1, &c2)
              if err9 != nil {
                      return err9
              }

              //fmt.Printf("%v  %v \n", c1, c2)
        }

        db.Query("DROP table " + tableOne)
        return nil
}

