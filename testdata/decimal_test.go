package main

import (
        "testing"
)

// Issue 116
func TestDecimalColumn(t *testing.T) {
	if DecimalColumn() != 1 {
		t.Error("Error at DecimalColumn")
	}
}
