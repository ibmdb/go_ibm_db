package main

import "testing"

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
                return err
        }
        return nil
}

