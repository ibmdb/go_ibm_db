package main

import (
	"testing"
)

func TestNullValueString(t *testing.T) {
	if NullValueString() != nil {
		t.Error("Error at NullValueString")
	}
}
