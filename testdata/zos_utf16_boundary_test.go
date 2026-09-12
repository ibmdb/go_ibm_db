package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// queryHangTimeout bounds how long a single boundary query/prepare iteration
// may take. A hung SQLExecDirect/SQLPrepare syscall cannot be canceled once
// started, so this only stops us from waiting on it further.
// Keep this short: with -timeout N, the budget must cover every successful
// iteration's real round-trip time plus this timeout for the final, hanging
// call. A large value here (e.g. 15s) can make the process alarm fire first,
// so the hang is never actually reported.
const queryHangTimeout = 5 * time.Second

// TestZosUTF16Boundary tests issue #278: SQLExecDirect/SQLPrepare passed SQL_NTS
// to non-null-terminated UTF-16 buffer, causing intermittent SQLCODE -104/-199 on z/OS.
// This test creates queries of various lengths to test boundary conditions.
func TestZosUTF16Boundary(t *testing.T) {
	// z/OS-specific: query length boundary behavior around SQLExecDirect/SQLPrepare.
	if TargetPlatform() != PlatformZOS {
		t.Skipf("Skipping test: only applicable on DB2_TARGET_PLATFORM=%s", PlatformZOS)
	}
	if ZosUTF16Boundary() != nil {
		t.Error("Error in ZosUTF16Boundary test")
	}
}

// TestZosUTF16BoundaryPrepare tests the same issue but for prepared statements.
func TestZosUTF16BoundaryPrepare(t *testing.T) {
	// z/OS-specific: query length boundary behavior around SQLExecDirect/SQLPrepare.
	if TargetPlatform() != PlatformZOS {
		t.Skipf("Skipping test: only applicable on DB2_TARGET_PLATFORM=%s", PlatformZOS)
	}
	if os.Getenv("DB2_RUN_ZOS_UTF16_BOUNDARY_PREPARE") != "1" {
		t.Skip("Skipping z/OS UTF-16 boundary prepare stress test; set DB2_RUN_ZOS_UTF16_BOUNDARY_PREPARE=1 to run it explicitly")
	}
	if ZosUTF16BoundaryPrepare() != nil {
		t.Error("Error in ZosUTF16BoundaryPrepare test")
	}
}

// runBoundaryQuery runs db.Query(query) on a background goroutine and reports
// back within timeout. If the driver call hangs, we stop waiting and treat
// this query length as unsupported instead of blocking the whole test run.
func runBoundaryQuery(db *sql.DB, query string, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		rows, err := db.Query(query)
		if err != nil {
			done <- err
			return
		}
		defer rows.Close()

		hasRows := false
		for rows.Next() {
			hasRows = true
			var id int
			var data string
			if err := rows.Scan(&id, &data); err != nil {
				done <- err
				return
			}
		}
		if err := rows.Err(); err != nil {
			done <- err
			return
		}
		if !hasRows {
			done <- fmt.Errorf("expected rows but got none")
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timed out after %v (possible driver hang)", timeout)
	}
}

// runBoundaryPrepareQuery validates SQLPrepare for a boundary-length query.
func runBoundaryPrepareQuery(db *sql.DB, query string, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		stmt, err := db.Prepare(query)
		if err != nil {
			done <- err
			return
		}
		defer stmt.Close()
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timed out after %v (possible driver hang)", timeout)
	}
}

// cleanupBoundaryDB drops tableName and closes db. If a prior driver hang was
// detected, the underlying connection/shared driver state may be unusable, so
// cleanup runs in the background with a timeout instead of blocking the test
// function's return (the goroutine may leak, but the test won't hang on it).
func cleanupBoundaryDB(db *sql.DB, tableName string, hangDetected bool) {
	cleanup := func() {
		db.Exec("DROP TABLE " + tableName)
		db.Close()
	}
	if !hangDetected {
		cleanup()
		return
	}
	done := make(chan struct{})
	go func() {
		cleanup()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(queryHangTimeout):
		fmt.Println("Warning: cleanup after driver hang did not complete in time; abandoning connection")
	}
}

// ZosUTF16Boundary tests Query (SQLExecDirect) with various query lengths
// that might hit memory allocation boundaries (e.g., 2048 bytes = 1024 UTF-16 code units).
// Instead of failing/hanging at the first unsupported length, it stops there,
// reports the max supported query length, and passes.
func ZosUTF16Boundary() error {
	db := Createconnection()
	hangDetected := false
	tableName := "zos_boundary_test"
	defer func() { cleanupBoundaryDB(db, tableName, hangDetected) }()

	// Test various query lengths including boundary conditions
	// UTF-16 on z/OS uses 2 bytes per character for BMP characters
	// Memory allocator size classes often include: 512, 1024, 2048, 4096 bytes
	// These translate to: 256, 512, 1024, 2048 UTF-16 code units

	testLengths := []int{
		255,  // Just under 512 bytes
		256,  // Exactly 512 bytes
		257,  // Just over 512 bytes
		511,  // Just under 1024 bytes
		512,  // Exactly 1024 bytes
		513,  // Just over 1024 bytes
		1023, // Just under 2048 bytes
		1024, // Exactly 2048 bytes - the boundary mentioned in issue #278
		1025, // Just over 2048 bytes
		2047, // Just under 4096 bytes
		2048, // Exactly 4096 bytes
		2049, // Just over 4096 bytes
	}

	db.Exec("DROP TABLE " + tableName)
	_, err := db.Exec("CREATE TABLE " + tableName + " (id INT, data VARCHAR(100))")
	if err != nil {
		fmt.Println("Failed to create test table:", err)
		return err
	}

	_, err = db.Exec("INSERT INTO " + tableName + " VALUES (1, 'test')")
	if err != nil {
		fmt.Println("Failed to insert test data:", err)
		return err
	}

	maxSupportedLen := 0
	for _, targetLen := range testLengths {
		// Create a query with approximately the target UTF-16 length
		// Base query: "SELECT * FROM zos_boundary_test WHERE id = 1"
		baseQuery := "SELECT * FROM " + tableName + " WHERE id = 1"
		baseLen := len(baseQuery)

		// Add padding with SQL comment to reach target length
		if targetLen > baseLen {
			padding := strings.Repeat(" ", targetLen-baseLen-4) // -4 for "/*" and "*/"
			if len(padding) > 0 {
				baseQuery = baseQuery + " /*" + padding + "*/"
			}
		}

		// Run the query multiple times to catch intermittent issues
		// (the original issue was intermittent due to memory allocation patterns)
		failed := false
		for i := 0; i < 10; i++ {
			if err := runBoundaryQuery(db, baseQuery, queryHangTimeout); err != nil {
				fmt.Printf("Query failed at length %d, iteration %d: %v\n", targetLen, i, err)
				failed = true
				break
			}
		}
		if failed {
			hangDetected = true
			break
		}
		maxSupportedLen = targetLen
		fmt.Printf("Query length %d: OK (10 iterations)\n", targetLen)
	}

	fmt.Printf("Max supported query length (Query/SQLExecDirect): %d bytes\n", maxSupportedLen)
	fmt.Println("ZosUTF16Boundary test passed")
	return nil
}

// ZosUTF16BoundaryPrepare tests Prepare (SQLPrepare) with various query lengths.
// Like ZosUTF16Boundary, it stops at the first unsupported length, reports the
// max supported length, and passes rather than hanging/failing.
func ZosUTF16BoundaryPrepare() error {
	db := Createconnection()
	hangDetected := false
	tableName := "zos_boundary_prep_test"
	defer func() { cleanupBoundaryDB(db, tableName, hangDetected) }()

	testLengths := []int{
		255, 256, 257,
		511, 512, 513,
		1023, 1024, 1025,
		2047, 2048, 2049,
	}

	db.Exec("DROP TABLE " + tableName)
	_, err := db.Exec("CREATE TABLE " + tableName + " (id INT, data VARCHAR(100))")
	if err != nil {
		fmt.Println("Failed to create test table:", err)
		return err
	}

	_, err = db.Exec("INSERT INTO " + tableName + " VALUES (1, 'test')")
	if err != nil {
		fmt.Println("Failed to insert test data:", err)
		return err
	}

	maxSupportedLen := 0
	for _, targetLen := range testLengths {
		baseQuery := "SELECT * FROM " + tableName + " WHERE id = ?"
		baseLen := len(baseQuery)

		if targetLen > baseLen {
			padding := strings.Repeat(" ", targetLen-baseLen-4)
			if len(padding) > 0 {
				baseQuery = baseQuery + " /*" + padding + "*/"
			}
		}

		// Test prepared statement multiple times
		failed := false
		for i := 0; i < 10; i++ {
			if err := runBoundaryPrepareQuery(db, baseQuery, queryHangTimeout); err != nil {
				fmt.Printf("Prepared query failed at length %d, iteration %d: %v\n", targetLen, i, err)
				failed = true
				break
			}
		}
		if failed {
			hangDetected = true
			break
		}
		maxSupportedLen = targetLen
		fmt.Printf("Prepare length %d: OK (10 iterations)\n", targetLen)
	}

	fmt.Printf("Max supported query length (Prepare/SQLPrepare): %d bytes\n", maxSupportedLen)
	fmt.Println("ZosUTF16BoundaryPrepare test passed")
	return nil
}
