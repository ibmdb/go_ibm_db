// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go_ibm_db

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/ibmdb/go_ibm_db/api"
	trc "github.com/ibmdb/go_ibm_db/log2"
)

type Rows struct {
	os *ODBCStmt
}

func (r *Rows) Columns() []string {
	trc.Trace1("rows.go: Columns()")

	names := make([]string, len(r.os.Cols))
	for i := 0; i < len(names); i++ {
		names[i] = r.os.Cols[i].Name()
	}
	return names
}

func (r *Rows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	//TODO(Akhil):This functions retuns the precision and scale of column.
	trc.Trace1("rows.go: ColumnTypePrecisionScale() - ENTRY")
	trc.Trace1(fmt.Sprintf("index=%d", index))

	ok = false;
	var namelen api.SQLSMALLINT
	namebuf := make([]byte, api.MAX_FIELD_SIZE)
	ret := api.SQLColAttribute(r.os.h, api.SQLUSMALLINT(index+1), api.SQL_DESC_TYPE_NAME, api.SQLPOINTER(unsafe.Pointer(&namebuf[0])), (api.MAX_FIELD_SIZE), (*api.SQLSMALLINT)(&namelen), (api.SQLPOINTER)(unsafe.Pointer(nil)))

	if IsError(ret) {
		fmt.Println(ret)
		trc.Trace1(fmt.Sprintf("Error: %s",ret))
		return 0, 0, false
	}
	dbtype := string(namebuf[:namelen])
	ret = api.SQLColAttribute(r.os.h, api.SQLUSMALLINT(index+1), api.SQL_DESC_PRECISION, api.SQLPOINTER(unsafe.Pointer(nil)), 0, (*api.SQLSMALLINT)(nil), (api.SQLPOINTER)(unsafe.Pointer(&precision)))
	if IsError(ret) {
		fmt.Println(ret)
		trc.Trace1(fmt.Sprintf("Error: %s",ret))
		return 0, 0, false
	}
	ret = api.SQLColAttribute(r.os.h, api.SQLUSMALLINT(index+1), api.SQL_DESC_SCALE, api.SQLPOINTER(unsafe.Pointer(nil)), 0, (*api.SQLSMALLINT)(nil), (api.SQLPOINTER)(unsafe.Pointer(&scale)))
	if IsError(ret) {
		fmt.Println(ret)
		trc.Trace1(fmt.Sprintf("Error: %s",ret))
		return 0, 0, false
	}
	if dbtype == "DECIMAL" {
		ok = true;
	} else if dbtype == "NUMERIC" {
		ok = true;
	} else if dbtype == "TIMESTAMP" {
		ok = true;
	}

	trc.Trace1(fmt.Sprintf("precision=%d, scale=%d", precision, scale))
	trc.Trace1("rows.go: ColumnTypePrecisionScale() - EXIT")
	return precision, scale, ok
}

func (r *Rows) ColumnTypeLength(index int) (length int64, ok bool) {
	//ToDo(Akhil):This functions retuns the length of column.
	trc.Trace1("rows.go: ColumnTypeLength() - ENTRY")

	ret := api.SQLColAttribute(r.os.h, api.SQLUSMALLINT(index+1), api.SQL_DESC_LENGTH, api.SQLPOINTER(unsafe.Pointer(nil)), 0, (*api.SQLSMALLINT)(nil), (api.SQLPOINTER)(unsafe.Pointer(&length)))
	if IsError(ret) {
		fmt.Println(ret)
		return 0, false
	}
	trc.Trace1("rows.go: ColumnTypeLength() - EXIT")
	return length, true
}

func (r *Rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	//TODO(Akhil):This functions retuns whether the column is nullable or not
	trc.Trace1("rows.go: ColumnTypeNullable() - ENTRY")

	var null int64
	ret := api.SQLColAttribute(r.os.h, api.SQLUSMALLINT(index+1), api.SQL_DESC_NULLABLE, api.SQLPOINTER(unsafe.Pointer(nil)), 0, (*api.SQLSMALLINT)(nil), (api.SQLPOINTER)(unsafe.Pointer(&null)))
	if IsError(ret) {
		fmt.Println(ret)
		return false, false
	}
	if null == api.SQL_NULLABLE {
		return true, true
	}
	trc.Trace1("rows.go: ColumnTypeNullable() - EXIT")
	return false, true
}

func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	//TODO(AKHIL):This function will return the scantype that can be used to scan
	//the data to the golang variable.
	trc.Trace1("rows.go: ColumnTypeScanType()")
	trc.Trace1(fmt.Sprintf("index=%d", index))

	a := r.os.Cols[index].TypeScan()
	return (a)
}

func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	//TODO(AKHIL):This functions retuns the dbtype(VARCHAR,DECIMAL etc..) of column.
	//namebuf can be of uint8 or byte
	trc.Trace1("rows.go: ColumnTypeDatabaseTypeName() - ENTRY")

	var namelen api.SQLSMALLINT
	namebuf := make([]byte, api.MAX_FIELD_SIZE)
	ret := api.SQLColAttribute(r.os.h, api.SQLUSMALLINT(index+1), api.SQL_DESC_TYPE_NAME, api.SQLPOINTER(unsafe.Pointer(&namebuf[0])), (api.MAX_FIELD_SIZE), (*api.SQLSMALLINT)(&namelen), (api.SQLPOINTER)(unsafe.Pointer(nil)))

	if IsError(ret) {
		fmt.Println(ret)
		return ""
	}
	dbtype := string(namebuf[:namelen])
	trc.Trace1("rows.go: ColumnTypeDatabaseTypeName() - EXIT")
	return dbtype
}

func (r *Rows) Next(dest []driver.Value) error {
	trc.Trace1("rows.go: Next() - ENTRY")

	ret := api.SQLFetch(r.os.h)
	if ret == api.SQL_NO_DATA {
		return io.EOF
	}
	if IsError(ret) {
		return NewError("SQLFetch", r.os.h)
	}
	for i := range dest {
		v, err := r.os.Cols[i].Value(r.os.h, i)
		if err != nil {
			return err
		}
		dest[i] = v
	}
	trc.Trace1("rows.go: Next() - EXIT")
	return nil
}

func (r *Rows) HasNextResultSet() bool {
	trc.Trace1("rows.go: HasVextResultSet()")

	return true
}

func (r *Rows) NextResultSet() error {
	trc.Trace1("rows.go: NextResultSet() - ENTRY")

	ret := api.SQLMoreResults(r.os.h)
	if ret == api.SQL_NO_DATA {
		return io.EOF
	}
	if IsError(ret) {
		return NewError("SQLMoreResults", r.os.h)
	}

	err := r.os.BindColumns()
	if err != nil {
		return err
	}
	trc.Trace1("rows.go: NextResultSet() - EXIT")
	return nil
}

func (r *Rows) Close() error {
	trc.Trace1("rows.go: Close()")

	return r.os.closeByRows()
}
