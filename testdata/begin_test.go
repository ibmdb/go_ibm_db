package main

import "testing"

func TestBegin(t *testing.T) {
	if Begin() != nil {
		t.Error("Error while displaying table")
	}
}
