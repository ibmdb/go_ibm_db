package main

import "testing"

func TestDecimal2Array(t *testing.T) {
	if Decimal2Array() != nil {
		t.Error("Error at Decimal2Array")
	}
}
