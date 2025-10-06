// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go_ibm_db

import (
        "errors"
        "database/sql/driver"
        "fmt"
        "strconv"
        trc "github.com/ibmdb/go_ibm_db/log2"
)

type Result struct {
        c *Conn
        rowCount int64
}

func (r *Result) LastInsertId() (int64, error) {
        trc.Trace1("result.go: LastInsetId() - ENTRY")

        lastInsertId := int64(0)

        s, err := r.c.Prepare("select identity_val_local() as last_inserted_id from sysibm.sysdummy1")
        if err != nil {
                return 0, err
        }
        defer s.Close()

        rows, err := s.Query(nil)
        if err != nil {
                return 0, err
        }
        defer rows.Close()

        dest := make([]driver.Value, 1)
        err = rows.Next(dest)
        if err != nil {
                return 0, err
        }

        for _, val := range dest {
           if v, ok := val.([]byte); ok {
               strVal := string(v)
               lastInsertId, err = strconv.ParseInt(strVal, 10, 64)
               if err != nil {
                  return 0, err
               }
           } else {
               return 0, errors.New("unexpected type")
           }
        }
        trc.Trace1(fmt.Sprintf("LastInsertid = %d", lastInsertId))
        trc.Trace1("result.go: LastInsetId() - EXIT")

        return lastInsertId, nil
}

func (r *Result) RowsAffected() (int64, error) {
	return r.rowCount, nil
}
