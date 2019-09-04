package main

import "testing"

func TestStoredProcedure(t *testing.T) {
	if StoredProcedure() != nil {
		t.Error("Error at stored procedure")
	}
}
