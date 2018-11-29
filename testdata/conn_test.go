package main

import "testing"

func TestCreateconnection(t *testing.T){
    if(Createconnection() == nil){
	t.Error("connection not formed")
}	
}