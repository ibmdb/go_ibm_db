package main

import (
        "testing"
)

func TestMultipleQuery(t *testing.T) {
	if MultipleQuery() != nil {
		t.Error("Error at MultipleQuery")
	}
}
