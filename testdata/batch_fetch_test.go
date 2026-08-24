package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func runBatchFetchTest(fetchSize int, totalRows int) error {
	var baseConnStr string
	var connStrFound bool
	baseConnStr, connStrFound = os.LookupEnv("DB2_CONNSTR")
	if !connStrFound {
		UpdateConnectionVariables()
		baseConnStr = "PROTOCOL=tcpip;HOSTNAME=" + host + ";PORT=" + port + ";DATABASE=" + database + ";UID=" + uid + ";PWD=" + pwd
	}
	connStrWithFetch := fmt.Sprintf("%s;FETCHSIZE=%d;", baseConnStr, fetchSize)

	db, err := sql.Open("go_ibm_db", connStrWithFetch)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %v", err)
	}
	defer db.Close()

	tableName := fmt.Sprintf("batch_test_%d", fetchSize)
	db.Exec("DROP TABLE " + tableName)
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s (id INT, name VARCHAR(50))", tableName))
	if err != nil {
		return fmt.Errorf("create table failed: %v", err)
	}
	defer db.Exec("DROP TABLE " + tableName)

	for i := 1; i <= totalRows; i++ {
		_, err = db.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%d, 'name_%d')", tableName, i, i))
		if err != nil {
			return fmt.Errorf("insert failed: %v", err)
		}
	}

	rows, err := db.Query(fmt.Sprintf("SELECT id, name FROM %s ORDER BY id", tableName))
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("scan failed at row %d: %v", count+1, err)
		}
		count++
		if id != count {
			return fmt.Errorf("expected id %d, got %d", count, id)
		}
		expectedName := fmt.Sprintf("name_%d", count)
		if name != expectedName {
			return fmt.Errorf("expected name %s, got %s", expectedName, name)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows.Err(): %v", err)
	}

	if count != totalRows {
		return fmt.Errorf("expected %d rows, got %d", totalRows, count)
	}

	return nil
}

func TestBatchFetch(t *testing.T) {
	testCases := []struct {
		name      string
		fetchSize int
		totalRows int
	}{
		{"FetchSize_1_SingleRow", 1, 10},
		{"FetchSize_10_ExactBatch", 10, 20},
		{"FetchSize_10_PartialBatch", 10, 25},
		{"FetchSize_100_LargeResultSet", 100, 250},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runBatchFetchTest(tc.fetchSize, tc.totalRows); err != nil {
				t.Errorf("TestBatchFetch failed for fetchSize=%d, totalRows=%d: %v", tc.fetchSize, tc.totalRows, err)
			}
		})
	}
}
