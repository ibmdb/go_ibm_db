package main

import (
        "testing"
)

func TestHugeQuery(t *testing.T) {
	if HugeQuery() != 1 {
		t.Error("Error at HugeQuery")
	}
}
