package main

import "testing"

func TestDropDB(t *testing.T) {
	if DropDB() != false {
		t.Error("Error while dropping Database")
	}
}
