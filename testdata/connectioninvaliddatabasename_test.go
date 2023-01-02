package main

import (
        "testing"
)

func TestConnectionInvalidDatabaseName(t *testing.T) {
	if ConnectionInvalidDatabaseName() != 1 {
		t.Error("Error at ConnectionInvalidDatabaseName")
	}
}
