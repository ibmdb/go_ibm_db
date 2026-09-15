package main

import (
	"fmt"
	"os"
	"testing"
)

// TestMain prints the detected DB2 platform once before any test case runs.
func TestMain(m *testing.M) {
	fmt.Printf("DB2 target platform: %s\n", TargetPlatform())
	os.Exit(m.Run())
}
