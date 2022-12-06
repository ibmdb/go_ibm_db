package main

import (
        "testing"
)

func TestBadConnectionString(t *testing.T) {
	if BadConnectionString() != 1 {
		t.Error("Error at BadConnectionString")
	}
}
