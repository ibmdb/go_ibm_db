package main

import "testing"

func TestDec2Array(t *testing.T) {
	if Dec2Array() != nil {
		t.Error("Error at Dec2Array")
	}
}
