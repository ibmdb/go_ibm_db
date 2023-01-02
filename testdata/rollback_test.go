package main

import "testing"

func TestRollbackTransaction(t *testing.T) {
	if RollbackTransaction() != nil {
		t.Error("Error in Rollback Transaction")
	}
}
