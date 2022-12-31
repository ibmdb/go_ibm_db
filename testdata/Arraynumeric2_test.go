package main

import "testing"

func TestNumeric2Array(t *testing.T) {
	if Numeric2Array() != nil {
		t.Error("Error at Numeric2Array")
	}
}
