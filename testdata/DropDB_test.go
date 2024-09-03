package main

import "testing"

func TestDropDB(t *testing.T) {
	if DropDB() != true {
		t.Error("Error while dropping Database")
	}
}
