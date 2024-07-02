// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go_ibm_db

import (
	"database/sql/driver"
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"github.com/ibmdb/go_ibm_db/api"
	trc "github.com/ibmdb/go_ibm_db/log2"
)

func IsError(ret api.SQLRETURN) bool {
    trc.Trace1("error.go: IsError() - ENTRY")
	if ret == api.SQL_SUCCESS {
	    trc.Trace1("api.SQL_SUCCESS")
	} else if ret == api.SQL_SUCCESS_WITH_INFO {
	    trc.Trace1("api.SQL_SUCCESS_WITH_INFO")
	}
	trc.Trace1("error.go: IsError() - EXIT")
	return !(ret == api.SQL_SUCCESS || ret == api.SQL_SUCCESS_WITH_INFO)
}

type DiagRecord struct {
	State       string
	NativeError int
	Message     string
}

func (r *DiagRecord) String() string {
	return fmt.Sprintf("{%s} %s", r.State, r.Message)
}

type Error struct {
	APIName string
	Diag    []DiagRecord
}

func (e *Error) Error() string {
    trc.Trace1("error.go: Error() - ENTRY")
	ss := make([]string, len(e.Diag))
	for i, r := range e.Diag {
		ss[i] = r.String()
	}
	trc.Trace1(fmt.Sprintf("%s : %s", e.APIName, ss))
	trc.Trace1("error.go: Error() - EXIT")
	return e.APIName + ": " + strings.Join(ss, "\n")
}

func NewError(apiName string, handle interface{}) error {
	trc.Trace1("error.go: NewError() - ENTRY")
	trc.Trace1(fmt.Sprintf("apiName=%s",apiName))
    
	var ret api.SQLRETURN
	h, ht := ToHandleAndType(handle)
	err := &Error{APIName: apiName}
	var ne api.SQLINTEGER
	state := make([]uint16, 6)
	msg := make([]uint16, api.SQL_MAX_MESSAGE_LENGTH)
	for i := 1; ; i++ {
		if runtime.GOOS == "zos" {
			ret = api.SQLGetDiagRec(ht, h, api.SQLSMALLINT(i),
				(*api.SQLWCHAR)(unsafe.Pointer(&state[0])), &ne,
				(*api.SQLWCHAR)(unsafe.Pointer(&msg[0])),
				api.SQLSMALLINT(2*len(msg)), nil) // odbc api on zos doesn't handle null terminated strings, the exact size is passed
		} else {
			ret = api.SQLGetDiagRec(ht, h, api.SQLSMALLINT(i),
				(*api.SQLWCHAR)(unsafe.Pointer(&state[0])), &ne,
				(*api.SQLWCHAR)(unsafe.Pointer(&msg[0])),
				api.SQLSMALLINT(len(msg)), nil)
		}
		if ret == api.SQL_NO_DATA {
			break
		}
		if IsError(ret) {
		    trc.Trace1(fmt.Sprintf("SQLGetDiagRec failed: ret=%d", ret))
			panic(fmt.Errorf("SQLGetDiagRec failed: ret=%d", ret))
		}
		r := DiagRecord{
			State:       api.UTF16ToString(state),
			NativeError: int(ne),
			Message:     api.UTF16ToString(msg),
		}
		if strings.Contains(r.Message, "CLI0106E") ||
			strings.Contains(r.Message, "CLI0107E") ||
			strings.Contains(r.Message, "CLI0108E") {
			return driver.ErrBadConn
		}
		err.Diag = append(err.Diag, r)
	}
	trc.Trace1(fmt.Sprintf("Error: %s", err))
	trc.Trace1("error.go: NewError() - EXIT")
	return err
}
