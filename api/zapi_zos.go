// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build zos

package api

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/ibmruntimes/go-recordio/v2/utils"
	trc "github.com/ibmdb/go_ibm_db/log2"
)

var dll utils.Dll

func init() {
	var e error
	e = dll.Open("DSNAO64C")
	if e != nil {
		fmt.Fprintf(os.Stderr, "Failed to load DLL %s\n", "DSNAO64C")
		os.Exit(1)
	}
}

func SQLAllocHandle(handleType SQLSMALLINT, inputHandle SQLHANDLE, outputHandle *SQLHANDLE) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLAllocHandle() - ENTRY")
	trc.Trace1(fmt.Sprintf("handleType = %d", handleType))

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLAllocHandle"), uintptr(handleType), uintptr(inputHandle), uintptr(unsafe.Pointer(outputHandle)))

	trc.Trace1(fmt.Sprintf("r = %d", r))
	trc.Trace1("api/zapi_zos.go SQLAllocHandle() - EXIT")
	return SQLRETURN(r)
}

func SQLBindCol(statementHandle SQLHSTMT, columnNumber SQLUSMALLINT, targetType SQLSMALLINT, targetValuePtr []byte, bufferLength SQLLEN, vallen *SQLLEN) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLBindCol() - ENTRY")
	trc.Trace1(fmt.Sprintf("columnNumber=%d, targetType=%d, bufferLength=%d, vallen=%d", columnNumber, targetType, bufferLength, vallen))

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLBindCol"), uintptr(statementHandle), uintptr(columnNumber), uintptr(targetType), uintptr(unsafe.Pointer(&targetValuePtr[0])), uintptr(bufferLength), uintptr(unsafe.Pointer(vallen)))

	trc.Trace1("api/zapi_zos.go SQLBindCol() - EXIT")
	return SQLRETURN(r)
}

func SQLBindParameter(statementHandle SQLHSTMT, parameterNumber SQLUSMALLINT, inputOutputType SQLSMALLINT, valueType SQLSMALLINT, parameterType SQLSMALLINT, columnSize SQLULEN, decimalDigits SQLSMALLINT, parameterValue SQLPOINTER, bufferLength SQLLEN, ind *SQLLEN) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLBindParameter() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLBindParameter"), uintptr(statementHandle), uintptr(parameterNumber), uintptr(inputOutputType), uintptr(valueType), uintptr(parameterType), uintptr(columnSize), uintptr(decimalDigits), uintptr(parameterValue), uintptr(bufferLength), uintptr(unsafe.Pointer(ind)))

	trc.Trace1("api/zapi_zos.go SQLBindParameter() - EXIT")
	return SQLRETURN(r)
}

func SQLCloseCursor(statementHandle SQLHSTMT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLCloseCursor() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLCloseCursor"), uintptr(statementHandle))

	trc.Trace1("api/zapi_zos.go SQLCloseCursor() - EXIT")
	return SQLRETURN(r)
}

func SQLDescribeCol(statementHandle SQLHSTMT, columnNumber SQLUSMALLINT, columnName *SQLWCHAR, bufferLength SQLSMALLINT, nameLengthPtr *SQLSMALLINT, dataTypePtr *SQLSMALLINT, columnSizePtr *SQLULEN, decimalDigitsPtr *SQLSMALLINT, nullablePtr *SQLSMALLINT) SQLRETURN {
	trc.Trace1("api/zapi_zos.go SQLDescribeCol() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDescribeColW"), uintptr(statementHandle), uintptr(columnNumber), uintptr(unsafe.Pointer(columnName)), uintptr(bufferLength), uintptr(unsafe.Pointer(nameLengthPtr)), uintptr(unsafe.Pointer(dataTypePtr)), uintptr(unsafe.Pointer(columnSizePtr)), uintptr(unsafe.Pointer(decimalDigitsPtr)), uintptr(unsafe.Pointer(nullablePtr)))

	trc.Trace1("api/zapi_zos.go SQLDescribeCol() - EXIT")
	return SQLRETURN(r)
}

func SQLDescribeParam(statementHandle SQLHSTMT, parameterNumber SQLUSMALLINT, dataTypePtr *SQLSMALLINT, parameterSizePtr *SQLULEN, decimalDigitsPtr *SQLSMALLINT, nullablePtr *SQLSMALLINT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLDescribeParam() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDescribeParam"), uintptr(statementHandle), uintptr(parameterNumber), uintptr(unsafe.Pointer(dataTypePtr)), uintptr(unsafe.Pointer(parameterSizePtr)), uintptr(unsafe.Pointer(decimalDigitsPtr)), uintptr(unsafe.Pointer(nullablePtr)))

	trc.Trace1("api/zapi_zos.go SQLDescribeParam() - EXIT")
	return SQLRETURN(r)
}

func SQLDisconnect(connectionHandle SQLHDBC) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLDisconnect() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDisconnect"), uintptr(connectionHandle))

	trc.Trace1("api/zapi_zos.go SQLDisconnect() - EXIT")
	return SQLRETURN(r)
}

func SQLDriverConnect(connectionHandle SQLHDBC, windowHandle SQLHWND, inConnectionString *SQLWCHAR, stringLength1 SQLSMALLINT, outConnectionString *SQLWCHAR, bufferLength SQLSMALLINT, stringLength2Ptr *SQLSMALLINT, driverCompletion SQLUSMALLINT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLDriverConnect() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDriverConnectW"),uintptr(connectionHandle), uintptr(windowHandle), uintptr(unsafe.Pointer(inConnectionString)), uintptr(stringLength1), uintptr(unsafe.Pointer(outConnectionString)), uintptr(bufferLength), uintptr(unsafe.Pointer(stringLength2Ptr)), uintptr(driverCompletion))

	trc.Trace1("api/zapi_zos.go SQLDriverConnect() - EXIT")
	return SQLRETURN(r)
}

func SQLEndTran(handleType SQLSMALLINT, handle SQLHANDLE, completionType SQLSMALLINT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLEndTran() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLEndTran"), uintptr(handleType), uintptr(handle), uintptr(completionType))

	trc.Trace1("api/zapi_zos.go SQLEndTran() - EXIT")
	return SQLRETURN(r)
}

func SQLExecute(statementHandle SQLHSTMT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLExecute() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLExecute"), uintptr(statementHandle))

	trc.Trace1("api/zapi_zos.go SQLExecute() - EXIT")
	return SQLRETURN(r)
}

func SQLFetch(statementHandle SQLHSTMT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLFetch() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLFetch"), uintptr(statementHandle))

	trc.Trace1("api/zapi_zos.go SQLFetch() - EXIT")
	return SQLRETURN(r)
}

func SQLFreeHandle(handleType SQLSMALLINT, handle SQLHANDLE) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLFreeHandle() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLFreeHandle"), uintptr(handleType), uintptr(handle))

	trc.Trace1("api/zapi_zos.go SQLFreeHandle() - EXIT")
	return SQLRETURN(r)
}

func SQLGetData(statementHandle SQLHSTMT, colOrParamNum SQLUSMALLINT, targetType SQLSMALLINT, targetValuePtr SQLPOINTER, bufferLength SQLLEN, vallen *SQLLEN) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLGetData() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLGetData"), uintptr(statementHandle), uintptr(colOrParamNum), uintptr(targetType), uintptr(targetValuePtr), uintptr(bufferLength), uintptr(unsafe.Pointer(vallen)))

	trc.Trace1("api/zapi_zos.go SQLGetData() - EXIT")
	return SQLRETURN(r)
}

func SQLGetDiagRec(handleType SQLSMALLINT, handle SQLHANDLE, recNumber SQLSMALLINT, sqlState *SQLWCHAR, nativeErrorPtr *SQLINTEGER, messageText *SQLWCHAR, bufferLength SQLSMALLINT, textLengthPtr *SQLSMALLINT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLGetDiagRec() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLGetDiagRecW"), uintptr(handleType), uintptr(handle), uintptr(recNumber), uintptr(unsafe.Pointer(sqlState)), uintptr(unsafe.Pointer(nativeErrorPtr)), uintptr(unsafe.Pointer(messageText)), uintptr(bufferLength), uintptr(unsafe.Pointer(textLengthPtr)))

	trc.Trace1("api/zapi_zos.go SQLGetDiagRec() - EXIT")
	return SQLRETURN(r)
}

func SQLNumParams(statementHandle SQLHSTMT, parameterCountPtr *SQLSMALLINT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLNumParams() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLNumParams"), uintptr(statementHandle), uintptr(unsafe.Pointer(parameterCountPtr)))

	trc.Trace1("api/zapi_zos.go SQLNumParams() - EXIT")
	return SQLRETURN(r)
}

func SQLNumResultCols(statementHandle SQLHSTMT, columnCountPtr *SQLSMALLINT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLNumResultCols() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLNumResultCols"), uintptr(statementHandle), uintptr(unsafe.Pointer(columnCountPtr)))

	trc.Trace1("api/zapi_zos.go SQLNumResultCols() - EXIT")
	return SQLRETURN(r)
}

func SQLPrepare(statementHandle SQLHSTMT, statementText *SQLWCHAR, textLength SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLPrepare() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLPrepareW"), uintptr(statementHandle), uintptr(unsafe.Pointer(statementText)), uintptr(textLength))

	trc.Trace1("api/zapi_zos.go SQLPrepare() - EXIT")
	return SQLRETURN(r)
}

func SQLRowCount(statementHandle SQLHSTMT, rowCountPtr *SQLLEN) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLRowCount() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLRowCount"), uintptr(statementHandle), uintptr(unsafe.Pointer(rowCountPtr)))

	trc.Trace1("api/zapi_zos.go SQLRowCount() - EXIT")
	return SQLRETURN(r)
}

func SQLSetEnvAttr(environmentHandle SQLHENV, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLSetEnvAttr() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetEnvAttr"), uintptr(environmentHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))

	trc.Trace1("api/zapi_zos.go SQLSetEnvAttr() - EXIT")
	return SQLRETURN(r)
}

func SQLSetConnectAttr(connectionHandle SQLHDBC, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLSetConnectAttr() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetConnectAttrW"), uintptr(connectionHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))

	trc.Trace1("api/zapi_zos.go SQLSetConnectAttr() - EXIT")
	return SQLRETURN(r)
}

func SQLColAttribute(statementHandle SQLHSTMT, ColumnNumber SQLUSMALLINT, FieldIdentifier SQLUSMALLINT, CharacterAttributePtr SQLPOINTER, BufferLength SQLSMALLINT, StringLengthPtr *SQLSMALLINT, NumericAttributePtr SQLPOINTER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLColAttribute() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLColAttribute"), uintptr(statementHandle), uintptr(ColumnNumber), uintptr(FieldIdentifier), uintptr(CharacterAttributePtr), uintptr(BufferLength), uintptr(unsafe.Pointer(StringLengthPtr)), uintptr(NumericAttributePtr))

	trc.Trace1("api/zapi_zos.go SQLColAttribute() - EXIT")
	return SQLRETURN(r)
}

func SQLMoreResults(statementHandle SQLHSTMT) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLMoreResults() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLMoreResults"), uintptr(statementHandle))

	trc.Trace1("api/zapi_zos.go SQLMoreResults() - EXIT")
	return SQLRETURN(r)
}

func SQLSetStmtAttr(statementHandle SQLHSTMT, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLSetStmtAttr() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetStmtAttr"), uintptr(statementHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))

	trc.Trace1("api/zapi_zos.go SQLSetStmtAttr() - EXIT")
	return SQLRETURN(r)
}

func SQLCreateDb(connectionHandle SQLHDBC, dbnamePtr *SQLWCHAR, dbnameLen SQLINTEGER, codeSetPtr *SQLWCHAR, codeSetLen SQLINTEGER, modePtr *SQLWCHAR, modeLen SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLCreateDb() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLCreateDbW"), uintptr(connectionHandle), uintptr(unsafe.Pointer(dbnamePtr)), uintptr(dbnameLen), uintptr(unsafe.Pointer(codeSetPtr)), uintptr(codeSetLen), uintptr(unsafe.Pointer(modePtr)), uintptr(modeLen))

	trc.Trace1("api/zapi_zos.go SQLCreateDb() - EXIT")
	return SQLRETURN(r)
}

func SQLDropDb(connectionHandle SQLHDBC, dbnamePtr *SQLWCHAR, dbnameLen SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLDropDb() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDropDbW"), uintptr(connectionHandle), uintptr(unsafe.Pointer(dbnamePtr)), uintptr(dbnameLen))

	trc.Trace1("api/zapi_zos.go SQLDropDb() - EXIT")
	return SQLRETURN(r)
}

func SQLExecDirect(statementHandle SQLHSTMT, statementText *SQLWCHAR, textLength SQLINTEGER) (SQLRETURN) {
	trc.Trace1("api/zapi_zos.go SQLExecDirect() - ENTRY")

	r := utils.CfuncEbcdic(getFunc(&dll, "SQLExecDirectW"), uintptr(statementHandle), uintptr(unsafe.Pointer(statementText)), uintptr(textLength))

	trc.Trace1("api/zapi_zos.go SQLExecDirect() - EXIT")
	return SQLRETURN(r)
}
