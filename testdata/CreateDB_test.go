package main

import "testing"

func TestCreateDB(t *testing.T) {
	if CreateDB() != false {
		t.Error("Error while creating Database")
	}
}
