package main

import "testing"

func TestRealArray(t *testing.T) {
	if RealArray() != nil {
		t.Error("Error at RealArray")
	}
}
