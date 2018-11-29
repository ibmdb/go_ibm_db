package main

import "testing"

func TestColumns(t *testing.T){
    if(Columns() != nil){
	t.Error("Error in displaying Query")
}	
}