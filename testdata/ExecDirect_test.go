package main

import "testing"

func TestExecDirect(t *testing.T) {
	if ExecDirect() != nil {
		t.Error("Error in ExecDirect")
	}
}
