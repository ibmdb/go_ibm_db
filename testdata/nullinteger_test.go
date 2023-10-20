package main

import (
        "database/sql"
	"time"
        "fmt"
        "strings"
        "testing"
)


func TestNullValueInteger(t *testing.T) {
	if NullValueInteger() != nil {
		t.Error("Error at NullValueInteger")
	}
}


//NullValueInteger function performs
func NullValueInteger() error {
        var out1, out2 sql.NullString
        var out3 sql.NullInt64
        var out4 sql.NullBool
        var out5 sql.NullFloat64
        var out6 sql.NullTime
        db := Createconnection()
        defer db.Close()
        db.Exec("Drop table arr")
        _, err := db.Exec("create table arr(var1 character, var2 varchar(30), var3 integer, var4 boolean, var5 double, var6 timestamp)")
        if err != nil {
                return err
        }
        c1 := "a"
        c2 := "test"
        c4 := true
        c5 := 1.234
        c6 := time.Now()
        st, err := db.Prepare("Insert into arr(var1,var2,var4,var5,var6) values(?,?,?,?,?)")
        defer st.Close()
        if err != nil {
                return err
        }
        _, err = st.Query(c1, c2, c4, c5, c6)
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting NullValueInteger")
                return err
        }
        st1, err := db.Prepare("select * from arr")
        defer st1.Close()
        if err != nil {
                return err
        }
        rows, err := st1.Query()
        if err != nil {
                return err
        }
        for rows.Next() {
                err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
                if err != nil {
                        fmt.Println("Error while retrieving NullValueInteger")
                        return err
                }
                if out1.String != c1 {
                        fmt.Println("Data is mismatched at NullValueInteger - character")
                        return fmt.Errorf("Wrong data retrieved")
                }
                if out2.String != c2 {
                        fmt.Println("Data is mismatched at NullValueInteger - String")
                        return fmt.Errorf("Wrong data retrieved")
                }
                if out4.Bool != c4 {
                        fmt.Println("Data is mismatched at NullValueInteger - Bool")
                        return fmt.Errorf("Wrong data retrieved")
                }
                if out5.Float64 != c5 {
                        fmt.Println("Data is mismatched at NullValueInteger - Float")
                        return fmt.Errorf("Wrong data retrieved")
                }
               if out6.Time.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") != c6.Format("2006-01-02 15:04:05 PM -07:00 Jan Mon MST") {
                        fmt.Println("Data is mismatched at NullValueInteger - Time")
                        return fmt.Errorf("Wrong data retrieved")
                }
        }
        return nil
}

