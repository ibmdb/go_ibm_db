package main

import (
	"testing"
)

func TestDecfloat(t *testing.T) {
	if Decfloat() != nil {
		t.Error("Error at Decfloat")
	}
}
