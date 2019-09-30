package main

import (
	"testing"
)

func TestNullValueCharacter(t *testing.T) {
	if NullValueCharacter() != nil {
		t.Error("Error at NullValueCharacter")
	}
}
