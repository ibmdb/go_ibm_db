// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build zos

package api

import (
	"log"
	"runtime"
	"unsafe"
	"github.com/ibmruntimes/go-recordio/v2/utils"
	trc "github.com/ibmdb/go_ibm_db/log2"
	"fmt"
)

func getFunc(dll *utils.Dll, str string) uintptr {
	fp, e := dll.Sym(str)
	if e != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[FATAL] %s [%s:%s:%d]", e, runtime.FuncForPC(pc).Name(), fn, line)
	}
	return fp
}

const (
	SQL_MAX_OPTION_STRING_LENGTH = 256
	SQL_OV_ODBC3                 = 3
	SQL_ATTR_ODBC_VERSION        = 200
	SQL_DRIVER_NOPROMPT          = 0
	SQL_HANDLE_ENV               = 1
	SQL_HANDLE_DBC               = 2
	SQL_HANDLE_STMT              = 3
	SQL_SUCCESS                  = 0
	SQL_SUCCESS_WITH_INFO        = 1
	SQL_INVALID_HANDLE           = -2
	SQL_NO_DATA                  = 100
	SQL_NO_TOTAL                 = -4
	SQL_NTS                      = -3
	SQL_MAX_MESSAGE_LENGTH       = 1024
	SQL_NULL_HANDLE              = 0
	SQL_NULL_HENV                = 0
	SQL_NULL_HDBC                = 0
	SQL_NULL_HSTMT               = 0
	SQL_PARAM_INPUT              = 1
	SQL_PARAM_OUTPUT             = 4
	SQL_PARAM_INPUT_OUTPUT       = 2
	SQL_NULL_DATA                = -1
	SQL_DATA_AT_EXEC             = -2
	SQL_CHAR                     = 1
	SQL_NUMERIC                  = 2
	SQL_DECIMAL                  = 3
	SQL_INTEGER                  = 4
	SQL_SMALLINT                 = 5
	SQL_FLOAT                    = 6
	SQL_REAL                     = 7
	SQL_DOUBLE                   = 8
	SQL_DATETIME                 = 9
	SQL_DATE                     = 9
	SQL_TIME                     = 10
	SQL_VARCHAR                  = 12
	SQL_TYPE_DATE                = 91
	SQL_TYPE_TIME                = 92
	SQL_TYPE_TIMESTAMP           = 93
	SQL_TIMESTAMP                = 11
	SQL_LONGVARCHAR              = -1
	SQL_BINARY                   = -2
	SQL_VARBINARY                = -3
	SQL_LONGVARBINARY            = -4
	SQL_BIGINT                   = -5
	SQL_TINYINT                  = -6
	SQL_BIT                      = -7
	SQL_WCHAR                    = -8
	SQL_WVARCHAR                 = -9
	SQL_WLONGVARCHAR             = -10
	SQL_BLOB                     = -98
	SQL_CLOB                     = -99
	SQL_SIGNED_OFFSET            = -20
	SQL_UNSIGNED_OFFSET          = -22
	SQL_DBCLOB                   = -350
	SQL_XML                      = -370
	SQL_COMMIT                   = 0
	SQL_ROLLBACK                 = 1
	SQL_AUTOCOMMIT               = 102
	SQL_ATTR_AUTOCOMMIT          = 102
	SQL_AUTOCOMMIT_OFF           = 0
	SQL_AUTOCOMMIT_ON            = 1
	SQL_AUTOCOMMIT_DEFAULT       = 1
	SQL_DESC_PRECISION           = 1005
	SQL_DESC_SCALE               = 1006
	SQL_DESC_LENGTH              = 1003
	SQL_DESC_CONCISE_TYPE        = 2
	SQL_DESC_TYPE_NAME           = 14
	SQL_COLUMN_TYPE              = 2
	SQL_COLUMN_TYPE_NAME         = 14
	SQL_DESC_NULLABLE            = 1008
	SQL_NULLABLE                 = 1
	SQL_NO_NULLS                 = 0
	SQL_DECFLOAT                 = -360
	SQL_ATTR_PARAMSET_SIZE       = 22
	SQL_IS_UINTEGER              = -5
	SQL_IS_INTEGER               = -6
	SQL_ATTR_CONNECTION_POOLING  = 201
	SQL_ATTR_CP_MATCH            = 202
	SQL_CP_OFF                   = 0
	SQL_CP_ONE_PER_DRIVER        = 1
	SQL_CP_ONE_PER_HENV          = 2
	SQL_CP_DEFAULT               = 0
	SQL_CP_STRICT_MATCH          = 0
	SQL_CP_RELAXED_MATCH         = 1
	SQL_C_CHAR                   = 1
	SQL_C_LONG                   = 4
	SQL_C_SHORT                  = 5
	SQL_C_FLOAT                  = 7
	SQL_C_DOUBLE                 = 8
	SQL_C_NUMERIC                = 2
	SQL_C_DATE                   = 9
	SQL_C_TIME                   = 10
	SQL_C_TYPE_TIMESTAMP         = 93
	SQL_C_TIMESTAMP              = 11
	SQL_C_BINARY                 = -2
	SQL_C_BIT                    = -7
	SQL_C_WCHAR                  = -8
	SQL_C_DEFAULT                = 99
	SQL_C_SBIGINT                = -25
	SQL_C_UBIGINT                = -27
	SQL_C_DBCHAR                 = -350
	SQL_C_TYPE_DATE              = 91
	SQL_C_TYPE_TIME              = 92

	// TODO(hyder): Not defined in sqlext.h. Using windows value, but it is not supported.
	SQL_SS_XML = -152

	MAX_FIELD_SIZE = 1024
	SQL_BOOLEAN    = 16
)

type (
	SQLCHAR       byte
	SQLSCHAR      byte
	SQLINTEGER    int32
	SDWORD        int32
	SQLSMALLINT   int16
	SQLDOUBLE     float64
	SQLREAL       float32
	SQLBIGINT     int64
	ODBCINT64     int64
	SQLRETURN     int16
	UDWORD        uint32
	SQLUINTEGER   uint32
	SQLUSMALLINT  uint16
	SQLSETPOSIROW uint16
	UWORD         uint16
	SQLPOINTER    unsafe.Pointer
	SQLLEN        SQLINTEGER
	SQLULEN       SQLUINTEGER
	SQLDBCHAR     uint16
	SQLWCHAR      uint16
	SQLHANDLE     SQLINTEGER
	SQLHENV       SQLINTEGER
	SQLHDBC       SQLINTEGER
	SQLHSTMT      SQLINTEGER
	SQLHWND       uintptr
)

func SQLSetEnvUIntPtrAttr(environmentHandle SQLHENV, attribute SQLINTEGER, valuePtr uintptr, stringLength SQLINTEGER) (ret SQLRETURN) {
	trc.Trace1("api/api_zos.go SQLSetEnvUIntPtrAttr() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetEnvAttr"), uintptr(environmentHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))

	trc.Trace1(fmt.Sprintf("r = %d", r))
	trc.Trace1("api/api_zos.go SQLSetEnvUIntPtrAttr() - EXIT")
	return SQLRETURN(r)
}

func SQLSetConnectUIntPtrAttr(connectionHandle SQLHDBC, attribute SQLINTEGER, valuePtr uintptr, stringLength SQLINTEGER) (ret SQLRETURN) {
	trc.Trace1("api/api_zos.go SQLSetConnectUIntPtrAttr() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetConnectAttr"), uintptr(connectionHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))

	trc.Trace1(fmt.Sprintf("r = %d", r))
	trc.Trace1("api/api_zos.go SQLSetConnectUIntPtrAttr() - EXIT")
	return SQLRETURN(r)
}
