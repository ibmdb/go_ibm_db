package main

import "testing"

func TestGraphicArray(t *testing.T) {
	if GraphicArray() != nil {
		t.Error("Error at GraphicArray")
	}
}
