package main

import (
        "testing"
)

func TestConnectionInvalidUserPassword(t *testing.T) {
	if ConnectionInvalidUserPassword() != 1 {
		t.Error("Error at ConnectionInvalidUserPassword")
	}
}
