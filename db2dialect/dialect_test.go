// Package db2dialect provides Bun ORM support for IBM DB2
package db2dialect

import (
	"math"
	"testing"

	"github.com/uptrace/bun/dialect/feature"
)

// TestDialectConformance ensures Dialect can be used as a dialect
// We don't do a strict interface check since Bun dialects use specific types
func TestDialectConformance(t *testing.T) {
	d := New()

	// Verify key interface methods exist and work
	name := d.Name()
	if name == 0 {
		t.Fatal("Name() returned zero value")
	}
	if d.IdentQuote() == 0 {
		t.Fatal("IdentQuote() returned null byte")
	}
	if d.Features() == 0 {
		t.Fatal("Features() returned zero")
	}

	t.Logf("✓ Dialect implements expected interface methods (name: %v)", name)
}

// TestDialectName verifies the dialect reports its name correctly
func TestDialectName(t *testing.T) {
	d := New()
	name := d.Name()

	if name != db2Name {
		t.Errorf("Expected DB2 dialect name %d, got %d", db2Name, name)
	}
	t.Logf("✓ Dialect name: %v", name)
}

// TestIdentQuote verifies the correct identifier quote character
func TestIdentQuote(t *testing.T) {
	d := New()
	quote := d.IdentQuote()

	expected := byte('"')
	if quote != expected {
		t.Errorf("Expected quote %q, got %q", expected, quote)
	}
	t.Logf("✓ Identifier quote: %c", quote)
}

// TestFeaturesOffsetFetch verifies that OffsetFetch is enabled for DB2 pagination
func TestFeaturesOffsetFetch(t *testing.T) {
	d := New()
	features := d.Features()

	if !features.Has(feature.OffsetFetch) {
		t.Fatal("OffsetFetch feature not enabled - DB2 pagination will fail")
	}
	t.Log("✓ OffsetFetch feature enabled")
}

// TestFeaturesIdentity verifies that Identity is enabled so Bun emits
// GENERATED ... AS IDENTITY on autoincrement PKs and omits them from INSERTs.
func TestFeaturesIdentity(t *testing.T) {
	d := New()
	features := d.Features()

	if !features.Has(feature.Identity) {
		t.Fatal("Identity feature not enabled - autoincrement columns will not work on DB2")
	}
	t.Log("✓ Identity feature enabled")
}

// TestFeaturesComposite verifies all expected features are enabled
func TestFeaturesComposite(t *testing.T) {
	d := New()
	features := d.Features()

	expectedFeatures := []feature.Feature{
		feature.CTE,
		feature.WithValues,
		feature.SelectExists,
		feature.CompositeIn,
		feature.OffsetFetch,
		feature.Identity,
	}

	for _, f := range expectedFeatures {
		if !features.Has(f) {
			t.Errorf("Feature %v not enabled", f)
		}
	}
	t.Logf("✓ All expected features enabled")
}

// TestDefaultVarcharLen verifies the default VARCHAR length
func TestDefaultVarcharLen(t *testing.T) {
	d := New()
	length := d.DefaultVarcharLen()

	expected := 255
	if length != expected {
		t.Errorf("Expected default varchar length %d, got %d", expected, length)
	}
	t.Logf("✓ Default VARCHAR length: %d", length)
}

// TestDefaultSchema verifies the default schema name
func TestDefaultSchema(t *testing.T) {
	d := New()
	schema := d.DefaultSchema()

	expected := ""
	if schema != expected {
		t.Errorf("Expected default schema %q, got %q", expected, schema)
	}
	t.Logf("✓ Default schema: %s", schema)
}

// TestTargetOptions verifies target platform option configuration and catalog helper methods
func TestTargetOptions(t *testing.T) {
	// Default target (LUW)
	dLuw := New()
	if dLuw.Target() != TargetLUW {
		t.Errorf("Expected default target TargetLUW, got %v", dLuw.Target())
	}
	if dLuw.Target().String() != "LUW" {
		t.Errorf("Expected string 'LUW', got %q", dLuw.Target().String())
	}
	if dLuw.CatalogSchema() != "SYSCAT" {
		t.Errorf("Expected catalog schema 'SYSCAT' for LUW, got %q", dLuw.CatalogSchema())
	}

	// z/OS target
	dZos := New(WithTarget(TargetZOS))
	if dZos.Target() != TargetZOS {
		t.Errorf("Expected target TargetZOS, got %v", dZos.Target())
	}
	if dZos.Target().String() != "z/OS" {
		t.Errorf("Expected string 'z/OS', got %q", dZos.Target().String())
	}
	if dZos.CatalogSchema() != "SYSIBM" {
		t.Errorf("Expected catalog schema 'SYSIBM' for z/OS, got %q", dZos.CatalogSchema())
	}

	// IBM i target
	dIbmi := New(WithTarget(TargetIBMi))
	if dIbmi.Target() != TargetIBMi {
		t.Errorf("Expected target TargetIBMi, got %v", dIbmi.Target())
	}
	if dIbmi.Target().String() != "IBM i" {
		t.Errorf("Expected string 'IBM i', got %q", dIbmi.Target().String())
	}
	if dIbmi.CatalogSchema() != "QSYS2" {
		t.Errorf("Expected catalog schema 'QSYS2' for IBM i, got %q", dIbmi.CatalogSchema())
	}

	// Dummy table check
	if dZos.DummyTable() != "SYSIBM.SYSDUMMY1" {
		t.Errorf("Expected dummy table 'SYSIBM.SYSDUMMY1', got %q", dZos.DummyTable())
	}
}

// TestExplicitConstructors verifies platform-specific constructors avoid auto-detection.
func TestExplicitConstructors(t *testing.T) {
	tests := []struct {
		name   string
		new    func() *Dialect
		target TargetPlatform
	}{
		{"LUW", NewLUW, TargetLUW},
		{"z/OS", NewZOS, TargetZOS},
		{"IBM i", NewIBMi, TargetIBMi},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dialect := test.new()
			if dialect.Target() != test.target {
				t.Fatalf("expected target %v, got %v", test.target, dialect.Target())
			}
			if !dialect.targetSetExplicitly {
				t.Fatal("explicit constructor must disable target auto-detection")
			}
			if dialect.autoDetected {
				t.Fatal("explicit constructor should not mark target as auto-detected")
			}
		})
	}
}

// TestAppendOffsetLimit verifies pagination SQL generation
func TestAppendOffsetLimit(t *testing.T) {
	tests := []struct {
		name     string
		offset   int64
		limit    int64
		expected string
	}{
		{"no pagination", 0, 0, ""},
		{"limit only", 0, 10, " FETCH NEXT 10 ROWS ONLY"},
		{"offset only", 20, 0, " OFFSET 20 ROWS"},
		{"offset and limit", 20, 10, " OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := make([]byte, 0, 64)
			result := AppendOffsetLimit(b, test.offset, test.limit)
			resultStr := string(result)

			if resultStr != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, resultStr)
			}
		})
	}
	t.Log("✓ AppendOffsetLimit SQL generation verified")
}

// TestAppendOffsetLimitNegativeValuesPanic verifies negative offset/limit are rejected
func TestAppendOffsetLimitNegativeValuesPanic(t *testing.T) {
	tests := []struct {
		name   string
		offset int64
		limit  int64
	}{
		{"negative offset", -1, 10},
		{"negative limit", 20, -1},
		{"min int64 offset", math.MinInt64, 10},
		{"min int64 limit", 20, math.MinInt64},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for offset=%d, limit=%d", test.offset, test.limit)
				}
			}()
			b := make([]byte, 0, 64)
			AppendOffsetLimit(b, test.offset, test.limit)
		})
	}
}
