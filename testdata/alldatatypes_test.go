package main

import (
        "testing"
)

func TestAllDataTypes(t *testing.T) {
	if AllDataTypes() != nil {
		t.Error("Error at AllDataTypes")
	}
}
