package main

import (
        "testing"
)

func TestQueryCreateDropTable(t *testing.T) {
        if QueryCreateDropTable() != 1 {
                t.Error("Error at CreateDropTable")
        }
}
