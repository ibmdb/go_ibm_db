package main

import "testing"

func TestQueryContext(t *testing.T){
    if(QueryContext() != nil){
	   t.Error("table not displayed")
    }
}

//QueryContext will execute the prepared statement
func QueryContext() error {
        db := Createconnection()
        defer db.Close()
        st, _ := db.PrepareContext(ctx, "select * from rocket")
        _, err := st.QueryContext(ctx)
        if err != nil {
                return err
        }
        return nil
}

