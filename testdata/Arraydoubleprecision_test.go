package main

import "testing"

func TestDoubleprecisionArray(t *testing.T) {
	if DoubleprecisionArray() != nil {
		t.Error("Error at DoubleprecisionArray")
	}
}
