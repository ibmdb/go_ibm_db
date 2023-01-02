package main

import (
        "testing"
)

func TestConnectionInvalidPortNumber(t *testing.T) {
	if ConnectionInvalidPortNumber() != 1 {
		t.Error("Error at ConnectionInvalidPortNumber")
	}
}
