package main

import (
        "fmt"
        "testing"
)

func TestCreateDropTable(t *testing.T) {
        fmt.Println("TestCreateDropTable()")
        if CreateDropTable() != 1 {
                t.Error("Error at CreateDropTable")
        }
}

func CreateDropTable() int {

        db := Createconnection()

        _, err2 := db.Exec("DROP table VMSAMPLE")
        if err2 != nil {
               _, err3 := db.Exec("create table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
               if err3 != nil {
                         return 0
               }
        } else {
              _, err4 := db.Exec("create table VMSAMPLE(ID varchar(20),NAME varchar(20),LOCATION varchar(20),POSITION varchar(20))")
              if err4 != nil {
                        return 0
             }
       }
       fmt.Println("TABLE CREATED Successfully")

       st, err5 := db.Prepare("Insert into VMSAMPLE(ID,NAME,LOCATION,POSITION) values('3242','mike','hyd','manager')")
       if err5 != nil {
           return 0
       }
       st.Query()

        _, err6 := db.Exec("DROP table VMSAMPLE")
        if err6 != nil {
                return 0
        }
        fmt.Println("TABLE DROP Successfully")

        return 1
}

