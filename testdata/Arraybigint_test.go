package main

import (
	"testing"
)

func TestBigintArray(t *testing.T) {
	if BigintArray() != nil {
		t.Error("Error at BigintArray")
	}
}
