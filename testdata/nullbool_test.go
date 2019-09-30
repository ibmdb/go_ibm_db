package main

import (
	"testing"
)

func TestNullValueBool(t *testing.T) {
	if NullValueBool() != nil {
		t.Error("Error at NullValueBool")
	}
}
