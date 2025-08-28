package main

import (
	//"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestCodeunits32(t *testing.T) {
	if ChineseCharCodeunits32_1() != nil {
		t.Error("Error at ChineseCodeunits32")
	}
}

func ChineseCharCodeunits32_1() error {
	db := Createconnection()
	defer db.Close()
	db.Exec("Drop table TT")
	_, err := db.Exec("create table TT(C1 INTEGER NOT NULL, C2 VARCHAR(30 CODEUNITS32))")
	if err != nil {
		fmt.Println("ERROR: CREATE TABLE ")
		return err
	}
	st, err := db.Prepare("Insert into TT(C1, C2) values(1,'▒~@▒~L▒~I▒~[~[▒~T')")
	if err != nil {
		fmt.Println("ERROR: PREPARE TABLE ")
		return err
	}
	defer st.Close()
	_, err = st.Query()
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []string")
		fmt.Println("Error:- ", err)
		return err
	}
	/*
	   queryTypes := "select * from TT;"

	   rows, err := db.Query(queryTypes)
	   if err != nil {
	           fmt.Println("Error:", err)
	           return err
	   }

	   for rows.Next() {
	           var c1 int
	           var c2 string
	           err = rows.Scan(&c1, &c2)
	           if err != nil {
	                   fmt.Println("Error:", err)
	           }
	           fmt.Printf("C1=: %v\t C2: %v\n", c1, c2)
	   }

	*/
	return nil
}
