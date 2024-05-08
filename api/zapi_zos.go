// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build zos

package api

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/ibmruntimes/go-recordio/v2/utils"
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
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLAllocHandle"), uintptr(handleType), uintptr(inputHandle), uintptr(unsafe.Pointer(outputHandle)))
	return SQLRETURN(r)
}

func SQLBindCol(statementHandle SQLHSTMT, columnNumber SQLUSMALLINT, targetType SQLSMALLINT, targetValuePtr []byte, bufferLength SQLLEN, vallen *SQLLEN) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLBindCol"), uintptr(statementHandle), uintptr(columnNumber), uintptr(targetType), uintptr(unsafe.Pointer(&targetValuePtr[0])), uintptr(bufferLength), uintptr(unsafe.Pointer(vallen)))
	return SQLRETURN(r)
}

func SQLBindParameter(statementHandle SQLHSTMT, parameterNumber SQLUSMALLINT, inputOutputType SQLSMALLINT, valueType SQLSMALLINT, parameterType SQLSMALLINT, columnSize SQLULEN, decimalDigits SQLSMALLINT, parameterValue SQLPOINTER, bufferLength SQLLEN, ind *SQLLEN) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLBindParameter"), uintptr(statementHandle), uintptr(parameterNumber), uintptr(inputOutputType), uintptr(valueType), uintptr(parameterType), uintptr(columnSize), uintptr(decimalDigits), uintptr(parameterValue), uintptr(bufferLength), uintptr(unsafe.Pointer(ind)))
	return SQLRETURN(r)
}

func SQLCloseCursor(statementHandle SQLHSTMT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLCloseCursor"), uintptr(statementHandle))
	return SQLRETURN(r)
}

func SQLDescribeCol(statementHandle SQLHSTMT, columnNumber SQLUSMALLINT, columnName *SQLWCHAR, bufferLength SQLSMALLINT, nameLengthPtr *SQLSMALLINT, dataTypePtr *SQLSMALLINT, columnSizePtr *SQLULEN, decimalDigitsPtr *SQLSMALLINT, nullablePtr *SQLSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDescribeColW"), uintptr(statementHandle), uintptr(columnNumber), uintptr(unsafe.Pointer(columnName)), uintptr(bufferLength), uintptr(unsafe.Pointer(nameLengthPtr)), uintptr(unsafe.Pointer(dataTypePtr)), uintptr(unsafe.Pointer(columnSizePtr)), uintptr(unsafe.Pointer(decimalDigitsPtr)), uintptr(unsafe.Pointer(nullablePtr)))
	return SQLRETURN(r)
}

func SQLDescribeParam(statementHandle SQLHSTMT, parameterNumber SQLUSMALLINT, dataTypePtr *SQLSMALLINT, parameterSizePtr *SQLULEN, decimalDigitsPtr *SQLSMALLINT, nullablePtr *SQLSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDescribeParam"), uintptr(statementHandle), uintptr(parameterNumber), uintptr(unsafe.Pointer(dataTypePtr)), uintptr(unsafe.Pointer(parameterSizePtr)), uintptr(unsafe.Pointer(decimalDigitsPtr)), uintptr(unsafe.Pointer(nullablePtr)))
	return SQLRETURN(r)
}

func SQLDisconnect(connectionHandle SQLHDBC) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDisconnect"), uintptr(connectionHandle))
	return SQLRETURN(r)
}

func SQLDriverConnect(connectionHandle SQLHDBC, windowHandle SQLHWND, inConnectionString *SQLWCHAR, stringLength1 SQLSMALLINT, outConnectionString *SQLWCHAR, bufferLength SQLSMALLINT, stringLength2Ptr *SQLSMALLINT, driverCompletion SQLUSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDriverConnectW"),uintptr(connectionHandle), uintptr(windowHandle), uintptr(unsafe.Pointer(inConnectionString)), uintptr(stringLength1), uintptr(unsafe.Pointer(outConnectionString)), uintptr(bufferLength), uintptr(unsafe.Pointer(stringLength2Ptr)), uintptr(driverCompletion))
	return SQLRETURN(r)
}

func SQLEndTran(handleType SQLSMALLINT, handle SQLHANDLE, completionType SQLSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLEndTran"), uintptr(handleType), uintptr(handle), uintptr(completionType))
	return SQLRETURN(r)
}

func SQLExecute(statementHandle SQLHSTMT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLExecute"), uintptr(statementHandle))
	return SQLRETURN(r)
}

func SQLFetch(statementHandle SQLHSTMT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLFetch"), uintptr(statementHandle))
	return SQLRETURN(r)
}

func SQLFreeHandle(handleType SQLSMALLINT, handle SQLHANDLE) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLFreeHandle"), uintptr(handleType), uintptr(handle))
	return SQLRETURN(r)
}

func SQLGetData(statementHandle SQLHSTMT, colOrParamNum SQLUSMALLINT, targetType SQLSMALLINT, targetValuePtr SQLPOINTER, bufferLength SQLLEN, vallen *SQLLEN) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLGetData"), uintptr(statementHandle), uintptr(colOrParamNum), uintptr(targetType), uintptr(targetValuePtr), uintptr(bufferLength), uintptr(unsafe.Pointer(vallen)))
	return SQLRETURN(r)
}

func SQLGetDiagRec(handleType SQLSMALLINT, handle SQLHANDLE, recNumber SQLSMALLINT, sqlState *SQLWCHAR, nativeErrorPtr *SQLINTEGER, messageText *SQLWCHAR, bufferLength SQLSMALLINT, textLengthPtr *SQLSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLGetDiagRecW"), uintptr(handleType), uintptr(handle), uintptr(recNumber), uintptr(unsafe.Pointer(sqlState)), uintptr(unsafe.Pointer(nativeErrorPtr)), uintptr(unsafe.Pointer(messageText)), uintptr(bufferLength), uintptr(unsafe.Pointer(textLengthPtr)))
	return SQLRETURN(r)
}

func SQLNumParams(statementHandle SQLHSTMT, parameterCountPtr *SQLSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLNumParams"), uintptr(statementHandle), uintptr(unsafe.Pointer(parameterCountPtr)))
	return SQLRETURN(r)
}

func SQLNumResultCols(statementHandle SQLHSTMT, columnCountPtr *SQLSMALLINT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLNumResultCols"), uintptr(statementHandle), uintptr(unsafe.Pointer(columnCountPtr)))
	return SQLRETURN(r)
}

func SQLPrepare(statementHandle SQLHSTMT, statementText *SQLWCHAR, textLength SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLPrepareW"), uintptr(statementHandle), uintptr(unsafe.Pointer(statementText)), uintptr(textLength))
	return SQLRETURN(r)
}

func SQLRowCount(statementHandle SQLHSTMT, rowCountPtr *SQLLEN) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLRowCount"), uintptr(statementHandle), uintptr(unsafe.Pointer(rowCountPtr)))
	return SQLRETURN(r)
}

func SQLSetEnvAttr(environmentHandle SQLHENV, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetEnvAttr"), uintptr(environmentHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))
	return SQLRETURN(r)
}

func SQLSetConnectAttr(connectionHandle SQLHDBC, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetConnectAttrW"), uintptr(connectionHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))
	return SQLRETURN(r)
}

func SQLColAttribute(statementHandle SQLHSTMT, ColumnNumber SQLUSMALLINT, FieldIdentifier SQLUSMALLINT, CharacterAttributePtr SQLPOINTER, BufferLength SQLSMALLINT, StringLengthPtr *SQLSMALLINT, NumericAttributePtr SQLPOINTER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLColAttribute"), uintptr(statementHandle), uintptr(ColumnNumber), uintptr(FieldIdentifier), uintptr(CharacterAttributePtr), uintptr(BufferLength), uintptr(unsafe.Pointer(StringLengthPtr)), uintptr(NumericAttributePtr))
	return SQLRETURN(r)
}

func SQLMoreResults(statementHandle SQLHSTMT) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLMoreResults"), uintptr(statementHandle))
	return SQLRETURN(r)
}

func SQLSetStmtAttr(statementHandle SQLHSTMT, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLSetStmtAttr"), uintptr(statementHandle), uintptr(attribute), uintptr(valuePtr), uintptr(stringLength))
	return SQLRETURN(r)
}

func SQLCreateDb(connectionHandle SQLHDBC, dbnamePtr *SQLWCHAR, dbnameLen SQLINTEGER, codeSetPtr *SQLWCHAR, codeSetLen SQLINTEGER, modePtr *SQLWCHAR, modeLen SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLCreateDbW"), uintptr(connectionHandle), uintptr(unsafe.Pointer(dbnamePtr)), uintptr(dbnameLen), uintptr(unsafe.Pointer(codeSetPtr)), uintptr(codeSetLen), uintptr(unsafe.Pointer(modePtr)), uintptr(modeLen))
	return SQLRETURN(r)
}

func SQLDropDb(connectionHandle SQLHDBC, dbnamePtr *SQLWCHAR, dbnameLen SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLDropDbW"), uintptr(connectionHandle), uintptr(unsafe.Pointer(dbnamePtr)), uintptr(dbnameLen))
	return SQLRETURN(r)
}

func SQLExecDirect(statementHandle SQLHSTMT, statementText *SQLWCHAR, textLength SQLINTEGER) (SQLRETURN) {
	r := utils.CfuncEbcdic(getFunc(&dll, "SQLExecDirectW"), uintptr(statementHandle), uintptr(unsafe.Pointer(statementText)), uintptr(textLength))
	return SQLRETURN(r)
}
