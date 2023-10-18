package main

import (
        "testing"
)


func TestQueryInsertSelect(t *testing.T) {
        if QueryInsertSelect() != 1 {
                t.Error("Error at QueryInsertSelect")
        }
}

func QueryInsertSelect() int {

        db := Createconnection()
        err := QueryCreateTable(db)
        if err != nil {
                return 0
        }

        err = QueryInsertRow(db)
        if err != nil {
              return 0
        }

        err = QueryDisplayTable(db)
        if err != nil {
               return 0
        }

        err = QueryDropTable(db)
        if err != nil {
                return 0
        }

        return 1
}

