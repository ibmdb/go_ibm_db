package main

import "testing"

func TestQueryContext(t *testing.T){
    if(QueryContext() != nil){
	   t.Error("table not displayed")
    }	
}
