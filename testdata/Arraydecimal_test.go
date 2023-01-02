package main

import "testing"

func TestDecimalArray(t *testing.T) {
	if DecimalArray() != nil {
		t.Error("Error at DecimalArray")
	}
}
