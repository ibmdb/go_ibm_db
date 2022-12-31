package main

import "testing"

func TestVarcharArray(t *testing.T) {
	if VarcharArray() != nil {
		t.Error("Error at VarcharArray")
	}
}
