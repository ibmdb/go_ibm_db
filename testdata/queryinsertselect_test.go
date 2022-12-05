package main

import (
        "testing"
)

func TestQueryInsertSelect(t *testing.T) {
        if QueryInsertSelect() != 1 {
                t.Error("Error at QueryInsertSelect")
        }
}
