package main

import "testing"

func TestStoredProcedureInOut(t *testing.T) {
	if StoredProcedureInOut() != nil {
		t.Error("Error at stored procedure")
	}
}
