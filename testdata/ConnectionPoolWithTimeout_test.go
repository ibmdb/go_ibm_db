package main

import (
	"testing"
)

func TestConnectionPoolWithTimeout(t *testing.T) {
	SkipOnPlatform(t, PlatformZOS)
	if ConnectionPoolWithTimeout() == 0 {
		t.Error("Error in Connection pool with timeout")
	}
}
