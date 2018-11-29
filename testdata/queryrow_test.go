package main

import "testing"

func TestQueryrow(t *testing.T){
    if(Queryrow() != nil){
	   t.Error("values not displayed")
}	
}