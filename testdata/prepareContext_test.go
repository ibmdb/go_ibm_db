package main

import (
	"testing"
	"fmt"
)

func TestPrepareContext(t *testing.T){
    if(PrepareContext() != nil){
	t.Error("Error in preparing PrepareContext")
    }
}

//PrepareContext will prepare the statement
func PrepareContext() error {
        db := Createconnection()
        defer db.Close()

        _, err := db.PrepareContext(ctx, "select * from rocket")
        if err != nil {
                fmt.Println("PrepareContext error: ", err)
                return err
        }

        return nil
}

