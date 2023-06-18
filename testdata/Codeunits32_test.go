package main

import (
        "testing"
)

func TestCodeunits32(t *testing.T) {
	if ChineseCharCodeunits32() != nil {
		t.Error("Error at ChineseCodeunits32")
	}
}
