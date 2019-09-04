package main

import "testing"

func TestCreatetable(t *testing.T) {
	if Createtable() != nil {
		t.Error("Error while creating table")
	}
}
