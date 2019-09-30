package main

import (
	"testing"
)

func TestTimeStampArray(t *testing.T) {
	if TimeStampArray() != nil {
		t.Error("Error at StringArray")
	}
}
