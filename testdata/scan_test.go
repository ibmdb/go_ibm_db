package main

import "testing"

func TestScan(t *testing.T){
    if(Scan() != nil){
	t.Error("Error in Scanning Query")
}	
}