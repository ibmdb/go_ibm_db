package main

import (
	"testing"
)

func TestSmallintArray(t *testing.T) {
	if SmallintArray() != nil {
		t.Error("Error at SmallintArray")
	}
}
