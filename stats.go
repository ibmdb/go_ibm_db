// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go_ibm_db

import (
	"fmt"
	"sync"

	"github.com/ibmdb/go_ibm_db/api"
	trc "github.com/ibmdb/go_ibm_db/log2"
)

type Stats struct {
	EnvCount  int
	ConnCount int
	StmtCount int
	mu        sync.Mutex
}

func (s *Stats) updateHandleCount(handleType api.SQLSMALLINT, change int) {
	trc.Trace1("stats.go: updateHandleCount() - ENTRY")
	trc.Trace1(fmt.Sprintf("change=%d", change))

	s.mu.Lock()
	defer s.mu.Unlock()
	switch handleType {
	case api.SQL_HANDLE_ENV:
		s.EnvCount += change
	case api.SQL_HANDLE_DBC:
		s.ConnCount += change
	case api.SQL_HANDLE_STMT:
		s.StmtCount += change
	default:
	    trc.Trace1(fmt.Sprintf("unexpected handle type %d", handleType))
		panic(fmt.Errorf("unexpected handle type %d", handleType))
	}

	trc.Trace1("stats.go: updateHandleCount() - EXIT")
}
