package main

import (
        "testing"
)

func TestChineseChar(t *testing.T) {
	if ChineseChar() != nil {
		t.Error("Error at ChineseChar")
	}
}
