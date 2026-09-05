package go_ibm_db

import (
	"sync"
	"testing"
	"time"
)

// TestPoolSize verifies size() sums used and available connections for a pool.
func TestPoolSize(t *testing.T) {
	p := Pconnect("PoolSize=3")

	p.usedPool["dsn1"] = []*DBP{{}, {}}
	p.availablePool["dsn2"] = []*DBP{{}}

	if got := p.size(); got != 3 {
		t.Errorf("expected size 3, got %d", got)
	}
}

// TestPoolCapacityIndependence verifies that filling one pool to capacity
// does not affect another independently created pool's capacity accounting.
func TestPoolCapacityIndependence(t *testing.T) {
	poolA := Pconnect("PoolSize=1")
	poolB := Pconnect("PoolSize=1")

	poolA.usedPool["dsn"] = []*DBP{{}}

	poolA.mu.Lock()
	aFull := poolA.size() >= poolA.poolSize
	poolA.mu.Unlock()
	if !aFull {
		t.Fatalf("expected poolA to be at capacity")
	}

	poolB.mu.Lock()
	bUnderCapacity := poolB.size() < poolB.poolSize
	poolB.mu.Unlock()
	if !bUnderCapacity {
		t.Errorf("poolB capacity should be unaffected by poolA reaching its limit")
	}
}

func TestPoolOpenDoesNotReturnConnectionTwice(t *testing.T) {
	p := Pconnect("PoolSize=32")
	if !p.Init(1, "dsn") {
		t.Fatal("expected pool initialization to succeed")
	}

	results := make(chan *DBP, 32)
	var wg sync.WaitGroup
	for i := 0; i < cap(results); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- p.Open("dsn")
		}()
	}
	wg.Wait()
	close(results)

	seen := make(map[*DBP]bool)
	for db := range results {
		if db == nil {
			continue
		}
		if seen[db] {
			t.Fatalf("pool returned connection %p more than once", db)
		}
		seen[db] = true
	}
	if len(seen) > p.poolSize {
		t.Fatalf("pool returned %d connections for capacity %d", len(seen), p.poolSize)
	}
}

func TestPoolLifetimeIsIndependent(t *testing.T) {
	poolA := Pconnect("PoolSize=1")
	poolB := Pconnect("PoolSize=1")

	poolA.SetConnMaxLifetime(1)

	poolA.mu.Lock()
	lifetimeA := poolA.maxLifetime
	poolA.mu.Unlock()
	poolB.mu.Lock()
	lifetimeB := poolB.maxLifetime
	poolB.mu.Unlock()

	if lifetimeA != time.Second {
		t.Fatalf("expected poolA lifetime of one second, got %s", lifetimeA)
	}
	if lifetimeB != time.Duration(defaultConnMaxLifetime)*time.Second {
		t.Fatalf("expected poolB to retain the default lifetime, got %s", lifetimeB)
	}
}
