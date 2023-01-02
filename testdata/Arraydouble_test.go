package main

import "testing"

func TestDoubleArray(t *testing.T) {
	if DoubleArray() != nil {
		t.Error("Error at DoubleArray")
	}
}
