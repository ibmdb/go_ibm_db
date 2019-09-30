package main

import (
	"testing"
)

func TestBoolArray(t *testing.T) {
	if BoolArray() != nil {
		t.Error("Error at BoolArray")
	}
}
