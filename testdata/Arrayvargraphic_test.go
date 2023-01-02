package main

import "testing"

func TestVargraphicArray(t *testing.T) {
	if VargraphicArray() != nil {
		t.Error("Error at VargraphicArray")
	}
}
