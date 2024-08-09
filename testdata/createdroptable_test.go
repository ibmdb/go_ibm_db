package main

import (
        "fmt"
        "testing"
)

func TestCreateDropTable(t *testing.T) {
        if CreateDropTable() != 1 {
                t.Error("Error at CreateDropTable")
        }
}

func CreateDropTable() int {
        db := Createconnection()
	defer db.Close()

        db.Exec("DROP table VMSAMPLE")

        _, errExec := db.Exec("create table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
        if errExec != nil {
                fmt.Println("Exec error: ", errExec)
                return 0
        }

       fmt.Println("TABLE CREATED Successfully")

       st, err5 := db.Prepare("Insert into VMSAMPLE(ID,NAME,LOCATION,POSITION) values('3242','mike','hyd','manager')")
       if err5 != nil {
               fmt.Println("Prepare error: ", err5)
               return 0
       }
       st.Query()

        _, err6 := db.Exec("DROP table VMSAMPLE")
        if err6 != nil {
                fmt.Println("Exec error: ", err6)
                return 0
        }
        fmt.Println("TABLE DROP Successfully")

        return 1
}

