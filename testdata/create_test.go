package main

import "testing"

func TestCreatetable(t *testing.T){
    if(Createtable() != nil){
	   t.Error("table not formed")
}	
}