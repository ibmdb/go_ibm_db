package main

import "testing"

func TestIntArray(t *testing.T) {
	if IntArray() != nil {
		t.Error("Error at IntArray")
	}
}
