package main

import (
        "fmt"
        "strings"
        "testing"
)

func TestMultipleQuery(t *testing.T) {
	if MultipleQuery() != nil {
		t.Error("Error at MultipleQuery")
	}
}

func MultipleQuery() error {
        db := Createconnection()
        defer db.Close()
        db.Exec("Drop table arr")
        _, err := db.Exec("create table arr(PID bigint, C1 varchar(255), C2 varchar(255), C3 varchar(255))")
        if err != nil {
                return err
        }
        _, err = db.Query("Insert into arr values('1', 'PersonA', 'LastNameA', 'QA')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('2', 'PersonB', 'LastNameB', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('3', 'PersonC', 'LastNameC', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('4', 'PersonD', 'LastNameD', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("Insert into arr values('5', 'PersonE', 'LastNameE', 'Dev')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }

        _, err = db.Query("UPDATE arr SET C3 = 'QA Intern' where C2 = 'LastNameD'")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while updating []string")
                return err
        } else {
                fmt.Println("Update statement successful")
        }

        _, err = db.Query("SELECT count(*) from arr where PID = 7")
        if err != nil {
                return err
        } else {
                fmt.Println("Select statement successful")
        }

        _, err = db.Query("SELECT * from arr where C3 = 'QA Intern'")
        if err != nil {
                return err
        } else {
                fmt.Println("Select statement successful")
        }


        _, err = db.Query("DELETE from arr where PID = 5")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while deleting []string")
                return err
        } else {
                fmt.Println("Delete statement successful")
        }

        _, err = db.Query("INSERT into arr values('6', 'PersonF', 'LastNameF', 'QA Lead')")
        if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
                fmt.Println("Error while inserting []string")
                return err
        }else {
                fmt.Println("Insert statement successful")
        }

        return nil
}

