package main

import (
	"testing"
)

func TestNullValueFloat(t *testing.T) {
	if NullValueFloat() != nil {
		t.Error("Error at NullValueFloat")
	}
}
