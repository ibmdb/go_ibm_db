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

func TestClassifyPlatform(t *testing.T) {
	tests := []struct {
		name string
		dbms string
		want string
	}{
		{"LUW", "DB2/LINUXX8664", PlatformLUW},
		{"z/OS DB2", "DB2", PlatformZOS},
		{"z/OS DSN", "DSN11015", PlatformZOS},
		{"IBM i", "AS/400", PlatformAS400},
		{"unknown defaults to LUW", "OTHER", PlatformLUW},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := classifyPlatform(test.dbms); got != test.want {
				t.Fatalf("classifyPlatform(%q) = %q, want %q", test.dbms, got, test.want)
			}
		})
	}
}
