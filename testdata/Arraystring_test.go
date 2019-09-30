package main

import (
	"testing"
)

func TestStringArray(t *testing.T) {
	if StringArray() != nil {
		t.Error("Error at StringArray")
	}
}
