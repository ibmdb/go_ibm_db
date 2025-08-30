package main

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestNullValueTime(t *testing.T) {
	if NullValueTime() != nil {
		t.Error("Error at NullValueTime")
	}
}

// NullValueTime function performs
func NullValueTime() error {
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
		fmt.Println("Exec error: ", err)
		return err
	}
	c1 := "a"
	c2 := "test"
	c3 := int64(10)
	c4 := true
	c5 := 1.234
	st, err := db.Prepare("Insert into arr(var1,var2,var3,var4,var5) values(?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(c1, c2, c3, c4, c5)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting NullValueString")
		return err
	}
	st1, err := db.Prepare("select * from arr")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st1.Close()
	rows, err := st1.Query()
	if err != nil {
		fmt.Println("Query error: ", err)
		return err
	}
	for rows.Next() {
		err := rows.Scan(&out1, &out2, &out3, &out4, &out5, &out6)
		if err != nil {
			fmt.Println("Error while retrieving NullValueString")
			return err
		}
		if out1.String != c1 {
			fmt.Println("Data is mismatched at NullValueFloat - Character")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out2.String != c2 {
			fmt.Println("Data is mismatched at NullValueFloat - string")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out3.Int64 != c3 {
			fmt.Println("Data is mismatched at NullValueFloat - Int64")
			return fmt.Errorf("Wrong data retrieved")
		}
		if out4.Bool != c4 {
			fmt.Println("Data is mismatched at NullValueFloat - Bool")
		}
		if out5.Float64 != c5 {
			fmt.Println("Data is mismatched at NullValueFloat - Float64")
			return fmt.Errorf("Wrong data retrieved")
		}
	}
	return nil
}
