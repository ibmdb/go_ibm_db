package main

import "testing"

func TestPrepare(t *testing.T){
    if(Prepare() != nil){
	t.Error("Error in preparing Query")
}	
}