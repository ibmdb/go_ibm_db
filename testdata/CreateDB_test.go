package main

import "testing"

func TestCreateDB(t *testing.T) {
	if CreateDB() != true {
		t.Error("Error while creating Database")
	}
}
