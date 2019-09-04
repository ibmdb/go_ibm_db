package main

import "testing"

func TestDrop(t *testing.T) {
	if Drop() != nil {
		t.Error("Error while droping table")
	}
}
