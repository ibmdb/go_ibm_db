package main

import "testing"

func TestNumericArray(t *testing.T) {
	if NumericArray() != nil {
		t.Error("Error at NumericArray")
	}
}
