package main

import "testing"

func TestCreateDB(t *testing.T) {
	SkipOnPlatform(t, PlatformZOS)
	if CreateDB() != true {
		t.Error("Error while creating Database")
	}
}
