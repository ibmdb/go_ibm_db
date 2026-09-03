package go_ibm_db

import "testing"

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
