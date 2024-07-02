// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go_ibm_db

import (
	"fmt"
	"github.com/ibmdb/go_ibm_db/api"
	trc "github.com/ibmdb/go_ibm_db/log2"
)

func ToHandleAndType(handle interface{}) (h api.SQLHANDLE, ht api.SQLSMALLINT) {
	trc.Trace1("handle.go: ToHandleAndType() - ENTRY")

	switch v := handle.(type) {
	case api.SQLHENV:
		if v == api.SQLHENV(api.SQL_NULL_HANDLE) {
			ht = 0
		} else {
			ht = api.SQL_HANDLE_ENV
		}
		h = api.SQLHANDLE(v)
	case api.SQLHDBC:
		ht = api.SQL_HANDLE_DBC
		h = api.SQLHANDLE(v)
	case api.SQLHSTMT:
		ht = api.SQL_HANDLE_STMT
		h = api.SQLHANDLE(v)
	default:
		panic(fmt.Errorf("unexpected handle type %T", v))
	}
	trc.Trace1("handle.go: ToHandleAndType() - EXIT")
	return h, ht
}

func releaseHandle(handle interface{}) error {
	trc.Trace1("handle.go: releaseHandle() - ENTRY")

	h, ht := ToHandleAndType(handle)
	ret := api.SQLFreeHandle(ht, h)
	if ret == api.SQL_INVALID_HANDLE {
		return fmt.Errorf("SQLFreeHandle(%d, %d) returns SQL_INVALID_HANDLE", ht, h)
	}
	if IsError(ret) {
		return NewError("SQLFreeHandle", handle)
	}
	drv.Stats.updateHandleCount(ht, -1)

	trc.Trace1("handle.go: releaseHandle() - EXIT")
	return nil
}
