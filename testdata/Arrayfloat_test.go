package main

import "testing"

func TestFloatArray(t *testing.T) {
	if FloatArray() != nil {
		t.Error("Error at FloatArray")
	}
}
