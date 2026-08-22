// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go_ibm_db

import (
	"strings"
	"testing"
)

func TestParseDSN(t *testing.T) {
	tests := []struct {
		name              string
		dsn               string
		expectedFetchSize int
		mustNotContain    []string
		mustContain       []string
	}{
		{
			name:              "Standard DSN without fetch size",
			dsn:               "HOSTNAME=localhost;PORT=50000;DATABASE=testdb;UID=db2inst1;PWD=password;",
			expectedFetchSize: 1,
			mustNotContain:    []string{"FETCHSIZE", "ROWARRAYSIZE"},
			mustContain:       []string{"HOSTNAME=localhost", "DATABASE=testdb"},
		},
		{
			name:              "DSN with FETCHSIZE",
			dsn:               "HOSTNAME=localhost;PORT=50000;DATABASE=testdb;UID=db2inst1;PWD=password;FETCHSIZE=1000;",
			expectedFetchSize: 1000,
			mustNotContain:    []string{"FETCHSIZE"},
			mustContain:       []string{"HOSTNAME=localhost", "DATABASE=testdb", "UID=db2inst1"},
		},
		{
			name:              "DSN with ROWARRAYSIZE lowercase",
			dsn:               "HOSTNAME=localhost;PORT=50000;DATABASE=testdb;UID=db2inst1;PWD=password;rowarraysize=500;",
			expectedFetchSize: 500,
			mustNotContain:    []string{"rowarraysize", "ROWARRAYSIZE"},
			mustContain:       []string{"HOSTNAME=localhost", "DATABASE=testdb"},
		},
		{
			name:              "DSN with mixed case fetchsize",
			dsn:               "DATABASE=testdb;FetchSize=250;",
			expectedFetchSize: 250,
			mustNotContain:    []string{"FetchSize", "fetchsize"},
			mustContain:       []string{"DATABASE=testdb"},
		},
		{
			name:              "DSN with invalid FETCHSIZE value",
			dsn:               "DATABASE=testdb;FETCHSIZE=invalid;",
			expectedFetchSize: 1,
			mustNotContain:    []string{"FETCHSIZE"},
			mustContain:       []string{"DATABASE=testdb"},
		},
		{
			name:              "DSN with zero FETCHSIZE value",
			dsn:               "DATABASE=testdb;FETCHSIZE=0;",
			expectedFetchSize: 1,
			mustNotContain:    []string{"FETCHSIZE"},
			mustContain:       []string{"DATABASE=testdb"},
		},
		{
			name:              "DSN with negative FETCHSIZE value",
			dsn:               "DATABASE=testdb;FETCHSIZE=-10;",
			expectedFetchSize: 1,
			mustNotContain:    []string{"FETCHSIZE"},
			mustContain:       []string{"DATABASE=testdb"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanDSN, fetchSize, err := parseDSN(tt.dsn)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fetchSize != tt.expectedFetchSize {
				t.Errorf("expected fetchSize %d, got %d", tt.expectedFetchSize, fetchSize)
			}
			for _, notExpected := range tt.mustNotContain {
				if strings.Contains(strings.ToUpper(cleanDSN), strings.ToUpper(notExpected)) {
					t.Errorf("cleanDSN %q should not contain %q", cleanDSN, notExpected)
				}
			}
			for _, expected := range tt.mustContain {
				if !strings.Contains(cleanDSN, expected) {
					t.Errorf("cleanDSN %q should contain %q", cleanDSN, expected)
				}
			}
		})
	}
}
