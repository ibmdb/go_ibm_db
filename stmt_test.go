package go_ibm_db

import (
	"context"
	"database/sql/driver"
	"errors"
	"os"
	"testing"
)

// TestQueryContextCancelAfterBindNoLeak checks that a context canceled after
// Exec and BindColumns succeeded does not leak the statement handle. The
// error used to be returned after usedByRows was set, so no Rows object was
// ever created to release the statement and Close could not free the handle.
func TestQueryContextCancelAfterBindNoLeak(t *testing.T) {
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

	st, err := conn.Prepare("SELECT 1 FROM SYSIBM.SYSDUMMY1")
	if err != nil {
		t.Fatal(err)
	}
	stmt := st.(*Stmt)
	stmtsAfterPrepare := liveStmtHandles()

	// The context is already canceled, so the final check in query1 fires
	// after Exec and BindColumns have succeeded.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := stmt.QueryContext(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("QueryContext error = %v, want context.Canceled", err)
	}

	// With usedByRows wrongly set, Close keeps the handle alive.
	if err := stmt.Close(); err != nil {
		t.Fatal(err)
	}

	if got, want := liveStmtHandles(), stmtsAfterPrepare-1; got != want {
		t.Fatalf("live statement handles = %d, want %d: the canceled query leaked the handle",
			got, want)
	}
}

// TestQueryContextCancelAfterBindReuse checks that cancellation closes the
// cursor so the same statement can execute another query without leaking.
func TestQueryContextCancelAfterBindReuse(t *testing.T) {
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

	st, err := conn.Prepare("SELECT 1 FROM SYSIBM.SYSDUMMY1")
	if err != nil {
		t.Fatal(err)
	}
	stmt := st.(*Stmt)
	defer stmt.Close()
	stmtsAfterPrepare := liveStmtHandles()

	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := stmt.QueryContext(ctx, nil); !errors.Is(err, context.Canceled) {
			t.Fatalf("QueryContext error = %v, want context.Canceled", err)
		}

		rows, err := stmt.QueryContext(context.Background(), nil)
		if err != nil {
			t.Fatalf("QueryContext after cancellation: %v", err)
		}
		values := make([]driver.Value, 1)
		err = rows.Next(values)
		closeErr := rows.Close()
		if err != nil {
			t.Fatal(err)
		}
		if closeErr != nil {
			t.Fatal(closeErr)
		}
		if values[0] != int64(1) {
			t.Fatalf("query result = %v, want 1", values[0])
		}
		if got := liveStmtHandles(); got != stmtsAfterPrepare {
			t.Fatalf("live statement handles = %d, want %d", got, stmtsAfterPrepare)
		}
	}
}

func liveStmtHandles() int {
	drv.Stats.mu.Lock()
	defer drv.Stats.mu.Unlock()
	return drv.Stats.StmtCount
}
