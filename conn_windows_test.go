package go_ibm_db

import (
	"os"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/ibmdb/go_ibm_db/api"
)

// TestGetInfoBufferLength compares GetInfo with a larger SQLGetInfoW buffer.
// SQL_KEYWORDS provides a long string without creating database objects.
func TestGetInfoBufferLength(t *testing.T) {
	dsn := os.Getenv("DB2_CONNSTR")
	if dsn == "" {
		t.Skip("set DB2_CONNSTR to run against a live database")
	}
	c, err := drv.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	conn := c.(*Conn)

	const sqlKeywords = 89 // ODBC SQL_KEYWORDS info type.
	buf := make([]uint16, 4096)
	var outLen api.SQLSMALLINT
	ret := api.SQLGetInfo(conn.h, sqlKeywords,
		api.SQLPOINTER(unsafe.Pointer(&buf[0])), api.SQLSMALLINT(2*len(buf)), &outLen)
	if IsError(ret) {
		t.Fatal(NewError("SQLGetInfo", conn.h))
	}
	if outLen <= 126 {
		t.Skipf("SQL_KEYWORDS is too short to exercise the buffer boundary: %d bytes", outLen)
	}

	// GetInfo has room for 127 UTF-16 code units and a terminator.
	// Passing 128 bytes instead of 256 truncates this value to 63 units.
	want := api.UTF16ToString(buf[:127])
	got, err := conn.GetInfo(sqlKeywords)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("GetInfo(SQL_KEYWORDS) = %q (%d UTF-16 units), want %q (%d UTF-16 units)",
			got, len(utf16.Encode([]rune(got))), want, len(utf16.Encode([]rune(want))))
	}
	t.Logf("GetInfo returned %d UTF-16 units matching SQLGetInfoW", len(utf16.Encode([]rune(got))))
}
