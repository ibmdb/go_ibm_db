package main

import "testing"

func TestCommit(t *testing.T){
    if(Commit() != nil){
	t.Error("Error in commit query")
}	
}