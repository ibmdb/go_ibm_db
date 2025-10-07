package main

import (
	"fmt"
	"testing"
)

func TestLastInsertId(t *testing.T) {
	if LastInsertId() != nil {
		t.Error("Error at LastInsertId")
	}
}

// BoolArray function performs inserting bool array.
func LastInsertId() error {
	db := Createconnection()
	defer db.Close()

	db.Exec("Drop table vmusers")
	_, err := db.Exec("create table vmusers(id integer generated always as identity(start with 1,increment by 1) not null, name varchar(30),Primary key(id))")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	sys, err := db.Exec("Insert into vmusers(name) values('Name1')")
	if err != nil {
		fmt.Println("Exec Insert error: ", err)
		return err
	}
	lastId, err := sys.LastInsertId()
        if err != nil {
		fmt.Println(" Last insert id error: ", err)
                return(err)
        }
        fmt.Println("LastInsertId: ", lastId)

	return nil
}
