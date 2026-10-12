package main

import "testing"

func TestDropDB(t *testing.T) {
	SkipOnPlatform(t, PlatformZOS)
	if DropDB() != true {
		t.Error("Error while dropping Database")
	}
}
