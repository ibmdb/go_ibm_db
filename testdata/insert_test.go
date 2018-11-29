package main

import "testing"

func TestInsert(t *testing.T){
    if(Insert() != nil){
	   t.Error("table not formed")
}	
}