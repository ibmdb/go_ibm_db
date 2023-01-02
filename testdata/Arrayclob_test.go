package main

import "testing"

func TestClobArray(t *testing.T) {
	if ClobArray() != nil {
		t.Error("Error at ClobArray")
	}
}
