package main

import "testing"

func TestCharacterArray(t *testing.T) {
	if CharacterArray() != nil {
		t.Error("Error at CharacterArray")
	}
}
