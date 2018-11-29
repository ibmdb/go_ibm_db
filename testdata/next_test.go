package main

import "testing"

func TestNext(t *testing.T){
    if(Next() != nil){
	t.Error("Error in Scanning Query")
}	
}