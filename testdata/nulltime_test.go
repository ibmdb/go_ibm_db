package main

import (
	"testing"
)

func TestNullValueTime(t *testing.T) {
	if NullValueTime() != nil {
		t.Error("Error at NullValueTime")
	}
}
