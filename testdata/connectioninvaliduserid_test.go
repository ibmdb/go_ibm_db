package main

import (
        "testing"
)

func TestConnectionInvalidUserID(t *testing.T) {
	if ConnectionInvalidUserID() != 1 {
		t.Error("Error at ConnectionInvalidUserID")
	}
}
