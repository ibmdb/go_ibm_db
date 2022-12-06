package main

import (
        "testing"
)

func TestCreateDropTable(t *testing.T) {
        if CreateDropTable() != 1 {
                t.Error("Error at CreateDropTable")
        }
}
