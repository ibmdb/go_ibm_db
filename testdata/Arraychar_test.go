package main

import "testing"

func TestCharArray(t *testing.T) {
	if CharArray() != nil {
		t.Error("Error at CharArray")
	}
}
