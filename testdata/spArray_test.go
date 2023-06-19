package main

import "testing"

func TestStoredProcedureArray(t *testing.T) {
	if StoredProcedureArray() != nil {
		t.Error("Error at stored procedure array")
	}
}
