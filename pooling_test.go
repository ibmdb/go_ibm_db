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

func TestPoolReleasePreventsNewConnections(t *testing.T) {
	p := Pconnect("PoolSize=1")
	p.Release()

	if got := p.Open("dsn"); got != nil {
		t.Fatal("expected Open to return nil after Release")
	}
	if p.Init(1, "dsn") {
		t.Fatal("expected Init to fail after Release")
	}
}

func TestPoolReuseResetsIdleTimer(t *testing.T) {
	const lifetime = time.Second

	p := Pconnect("PoolSize=1")
	p.maxLifetime = lifetime
	if !p.Init(1, "dsn") {
		t.Fatal("expected pool initialization to succeed")
	}
	t.Cleanup(p.Release)

	dbp := p.Open("dsn")
	if dbp == nil {
		t.Fatal("expected Open to return a connection")
	}
	dbp.Close()

	time.Sleep(800 * time.Millisecond)
	reused := p.Open("dsn")
	if reused != dbp {
		t.Fatal("expected Open to reuse the available connection")
	}
	reused.Close()

	time.Sleep(400 * time.Millisecond)
	p.mu.Lock()
	available := len(p.availablePool["dsn"])
	p.mu.Unlock()
	if available != 1 {
		t.Fatalf("stale idle timer removed reused connection: got %d available, want 1", available)
	}

	deadline := time.Now().Add(lifetime)
	for time.Now().Before(deadline) {
		p.mu.Lock()
		available = len(p.availablePool["dsn"])
		p.mu.Unlock()
		if available == 0 {
			p.mu.Lock()
			timer := reused.idleTimer
			p.mu.Unlock()
			if timer != nil {
				t.Fatal("expired connection kept a non-nil idleTimer reference")
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("connection did not expire after the reset idle lifetime")
}

// TestPoolStaleExpireKeepsActiveTimer verifies that a callback from an old
// generation neither removes the connection nor clears the reference to the
// newer active timer.
func TestPoolStaleExpireKeepsActiveTimer(t *testing.T) {
	p := Pconnect("PoolSize=1")
	p.maxLifetime = time.Hour
	if !p.Init(1, "dsn") {
		t.Fatal("expected pool initialization to succeed")
	}
	t.Cleanup(p.Release)

	dbp := p.Open("dsn")
	if dbp == nil {
		t.Fatal("expected Open to return a connection")
	}
	dbp.Close()

	p.mu.Lock()
	staleGen := dbp.idleGen - 1
	activeTimer := dbp.idleTimer
	p.mu.Unlock()
	if activeTimer == nil {
		t.Fatal("expected an active idle timer after Close")
	}

	dbp.expire(staleGen)

	p.mu.Lock()
	defer p.mu.Unlock()
	if dbp.idleTimer != activeTimer {
		t.Fatal("stale expire callback cleared the newer active timer reference")
	}
	if got := len(p.availablePool["dsn"]); got != 1 {
		t.Fatalf("stale expire callback removed the connection: got %d available, want 1", got)
	}
}
