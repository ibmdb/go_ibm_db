package main

import (
	"testing"
)

func TestNullValueInteger(t *testing.T) {
	if NullValueInteger() != nil {
		t.Error("Error at NullValueInteger")
	}
}
