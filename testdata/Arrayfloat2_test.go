package main

import "testing"

func TestFloat2Array(t *testing.T) {
	if Float2Array() != nil {
		t.Error("Error at Float2Array")
	}
}
