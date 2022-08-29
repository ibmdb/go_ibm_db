package main

import "testing"

func TestPrepareContext(t *testing.T){
    if(PrepareContext() != nil){
	t.Error("Error in preparing PrepareContext")
}	
}
