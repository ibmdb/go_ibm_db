package main

import "testing"

func TestClose(t *testing.T){
    if(Close() != nil){
	t.Error("Error in Scanning Query")
}	
}