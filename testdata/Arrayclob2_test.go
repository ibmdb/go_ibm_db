package main

import "testing"

func TestClob2Array(t *testing.T) {
	if Clob2Array() != nil {
		t.Error("Error at Clob2Array")
	}
}
