package main

import "testing"

func TestInt2Array(t *testing.T) {
	if Int2Array() != nil {
		t.Error("Error at Int2Array")
	}
}
