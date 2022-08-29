package main

import "testing"

func TestExecContext(t *testing.T){
    if(Createtable_ExecContext() != nil){
	t.Error("Error in preparing ExecContext")
}	
}
