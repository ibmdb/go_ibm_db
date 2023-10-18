package main

import (
        "testing"
)

func TestQueryCreateDropTable(t *testing.T) {
        if QueryCreateDropTable() != 1 {
                t.Error("Error at CreateDropTable")
        }
}

func QueryCreateDropTable() int {
	db := Createconnection()
        if db == nil {
                return 0
        }

        err := QueryCreateTable(db)
        if err != nil {
                return 0
        }

        err = QueryDropTable(db)
        if err != nil {
                return 0
        }

        return 1
}

