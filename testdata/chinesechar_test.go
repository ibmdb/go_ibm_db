package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestChineseChar(t *testing.T) {
	if ChineseChar() != nil {
		t.Error("Error at ChineseChar")
	}
}

func ChineseChar() error {
	db := Createconnection()
	defer db.Close()

	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(ID bigint, var2 varchar(30))")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}
	//st, err := db.Prepare("Insert into arr values('101','2019年2▒~V~R~V~RH')")
	st, err := db.Prepare("Insert into arr values('101',x'32303139E5B9B431E69C88E4')")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query()
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []string")
		return err
	}

	return nil
}
