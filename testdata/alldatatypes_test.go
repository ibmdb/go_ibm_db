package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestAllDataTypes(t *testing.T) {
	if AllDataTypes() != nil {
		t.Error("Error at AllDataTypes")
	}
}

func AllDataTypes() error {
	db := Createconnection()
	defer db.Close()

	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr (c1 int, c2 SMALLINT, c3 BIGINT, c4 INTEGER, c5 DECIMAL(4,2), c6 NUMERIC, c7 float, c8 double, c9 decfloat, c10 char(10), c11 varchar(10), c12 char for bit data, c13 clob(10),c14 dbclob(100), c15 date, c16 time, c17 timestamp, c18 blob(10), c19 boolean) ccsid unicode")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	st, err := db.Prepare("insert into arr values (1, 2, 9007199254741997, 1234, 67.98, 5689, 56.2390, 34567890, 45.234, 'Vijay', 'Raj', '\x50', 'test123456','▒~V~R~P~@▒~V~R~P~A▒~V~R~P~B▒~V~R~P~C▒~V~R~P~D▒~V~R~P~E▒~V~R~P~F','2015-09-10', '10:16:33', '2015-09-10 10:16:33.770139', BLOB(x'616263'), true)")
	if err != nil {
		fmt.Println("Close error: ", err)
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
