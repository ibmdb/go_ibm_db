package main

import "testing"

func TestQuery(t *testing.T){
    if(Query() != nil){
	   t.Error("table not displayed")
}	
}