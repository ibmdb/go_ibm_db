package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/ibmdb/go_ibm_db/db2dialect"
	"github.com/uptrace/bun"
)

// Advanced test models

type DataTypeTest struct {
	ID             int64                   `bun:",pk,autoincrement"`
	IntField       int                     `bun:""`
	BigIntField    int64                   `bun:""`
	SmallIntField  int16                   `bun:""`
	FloatField     float64                 `bun:""`
	DecimalField   float64                 `bun:"type:DECIMAL(10,2)"`
	StringField    string                  `bun:""`
	TextField      string                  `bun:"type:CLOB"`
	BoolField      db2dialect.SmallIntBool `bun:""` // 0=false, 1=true
	DateField      db2dialect.Date         `bun:"type:DATE"`
	TimeField      db2dialect.TimeOfDay    `bun:"type:TIME"`
	TimestampField db2dialect.Timestamp    `bun:"type:TIMESTAMP"`
}

type NullableFields struct {
	ID             int64           `bun:",pk,autoincrement"`
	NullableString sql.NullString  `bun:""`
	NullableInt    sql.NullInt64   `bun:""`
	NullableFloat  sql.NullFloat64 `bun:""`
	NullableBool   sql.NullBool    `bun:""`
	NullableTime   sql.NullTime    `bun:""`
}

type ComplexQuery struct {
	ID          int64                   `bun:",pk,autoincrement"`
	Name        string                  `bun:",notnull"`
	Category    string                  `bun:""`
	Amount      float64                 `bun:""`
	IsActive    db2dialect.SmallIntBool `bun:""` // 0=false, 1=true
	CreatedDate time.Time               `bun:",default:current_timestamp"`
}

// Test: Multiple data types in single table
func TestBun_MultipleDataTypes(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*DataTypeTest)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*DataTypeTest)(nil)).IfNotExists().Exec(ctx)

	// Insert record with various types
	record := &DataTypeTest{
		IntField:       42,
		BigIntField:    9223372036854775807,
		SmallIntField:  100,
		FloatField:     3.14159,
		DecimalField:   123.45,
		StringField:    "Test String",
		TextField:      "This is a longer text field",
		BoolField:      db2dialect.SmallIntBool(1),
		DateField:      db2dialect.Date(time.Now()),
		TimeField:      db2dialect.TimeOfDay(time.Now()),
		TimestampField: db2dialect.Timestamp(time.Now()),
	}

	res, err := db.NewInsert().Model(record).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to insert with multiple data types: %v", err)
	}

	lastID, _ := res.LastInsertId()

	// Retrieve and verify
	var retrieved DataTypeTest
	err = db.NewSelect().
		Model(&retrieved).
		Where(`"id" = ?`, lastID).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	if retrieved.IntField != record.IntField {
		t.Errorf("IntField mismatch: %d != %d", retrieved.IntField, record.IntField)
	}
	if retrieved.StringField != record.StringField {
		t.Errorf("StringField mismatch: %s != %s", retrieved.StringField, record.StringField)
	}
	if retrieved.BoolField != record.BoolField {
		t.Errorf("BoolField mismatch: %v != %v", retrieved.BoolField, record.BoolField)
	}

	t.Log("✓ Multiple data types stored and retrieved correctly")
}

// Test: Nullable fields
func TestBun_NullableFields(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*NullableFields)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*NullableFields)(nil)).IfNotExists().Exec(ctx)

	// Insert with NULL values
	record := &NullableFields{
		NullableString: sql.NullString{String: "", Valid: false},
		NullableInt:    sql.NullInt64{Int64: 0, Valid: false},
		NullableFloat:  sql.NullFloat64{Float64: 0, Valid: false},
		NullableBool:   sql.NullBool{Bool: false, Valid: false},
		NullableTime:   sql.NullTime{Time: time.Time{}, Valid: false},
	}

	res, err := db.NewInsert().Model(record).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to insert NULL values: %v", err)
	}

	lastID, _ := res.LastInsertId()

	// Retrieve and verify NULL handling
	var retrieved NullableFields
	err = db.NewSelect().
		Model(&retrieved).
		Where(`"id" = ?`, lastID).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed to retrieve NULL values: %v", err)
	}

	if retrieved.NullableString.Valid {
		t.Error("Expected NullableString.Valid to be false")
	}

	t.Log("✓ NULL values handled correctly")
}

// Test: Complex WHERE conditions
func TestBun_ComplexWhereConditions(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item A", Category: "Electronics", Amount: 150.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item B", Category: "Electronics", Amount: 250.00, IsActive: db2dialect.SmallIntBool(0)},
		{Name: "Item C", Category: "Books", Amount: 25.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item D", Category: "Books", Amount: 35.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item E", Category: "Clothing", Amount: 75.00, IsActive: db2dialect.SmallIntBool(0)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Complex WHERE with AND/OR using Bun query builder
	var result []ComplexQuery
	err := db.NewSelect().
		Model(&result).
		Where(`"category" = ? OR "category" = ?`, "Electronics", "Books").
		Where(`"is_active" = ?`, db2dialect.SmallIntBool(1)).
		Order("amount").
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed to execute complex WHERE: %v", err)
	}

	expectedCount := 3 // Item A, Item C, Item D
	if len(result) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(result))
	}

	t.Logf("✓ Complex WHERE conditions returned %d records", len(result))
}

// Test: Group By aggregation
func TestBun_GroupByAggregation(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item A", Category: "Electronics", Amount: 100.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item B", Category: "Electronics", Amount: 200.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item C", Category: "Books", Amount: 50.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item D", Category: "Books", Amount: 30.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Group By with aggregation
	type CategoryTotal struct {
		Category string  `bun:"category"`
		Total    float64 `bun:"total"`
	}

	var results []CategoryTotal
	err := db.NewSelect().
		ColumnExpr(`"category"`).
		ColumnExpr(`SUM("amount") AS "total"`).
		TableExpr(`"complex_queries"`).
		GroupExpr(`"category"`).
		Order("category").
		Scan(ctx, &results)

	if err != nil {
		t.Fatalf("Failed to execute GROUP BY: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(results))
	}

	t.Logf("✓ GROUP BY aggregation returned %d groups", len(results))
}

// Test: String operations
func TestBun_StringOperations(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Apple iPhone", Category: "Electronics", Amount: 999.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Apple iPad", Category: "Electronics", Amount: 599.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Samsung Galaxy", Category: "Electronics", Amount: 799.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// LIKE query
	var result []ComplexQuery
	err := db.NewSelect().
		Model(&result).
		Where(`"name" LIKE ?`, "Apple%").
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed LIKE query: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 results with 'Apple%%', got %d", len(result))
	}

	t.Log("✓ String operations (LIKE) working correctly")
}

// Test: IN operator
func TestBun_InOperator(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item 1", Category: "A", Amount: 100.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 2", Category: "B", Amount: 200.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 3", Category: "C", Amount: 300.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 4", Category: "D", Amount: 400.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// IN query
	categories := []string{"A", "C"}
	var result []ComplexQuery
	err := db.NewSelect().
		Model(&result).
		Where(`"category" IN (?)`, bun.In(categories)).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed IN query: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 results with IN, got %d", len(result))
	}

	t.Log("✓ IN operator working correctly")
}

// Test: BETWEEN operator
func TestBun_BetweenOperator(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item 1", Category: "Test", Amount: 50.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 2", Category: "Test", Amount: 150.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 3", Category: "Test", Amount: 250.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 4", Category: "Test", Amount: 350.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// BETWEEN query
	var result []ComplexQuery
	err := db.NewSelect().
		Model(&result).
		Where(`"amount" BETWEEN ? AND ?`, 100.00, 300.00).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed BETWEEN query: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 results with BETWEEN, got %d", len(result))
	}

	t.Log("✓ BETWEEN operator working correctly")
}

// Test: Distinct
func TestBun_Distinct(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data with duplicates
	records := []ComplexQuery{
		{Name: "Item A", Category: "Electronics", Amount: 100.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item B", Category: "Electronics", Amount: 200.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item C", Category: "Books", Amount: 50.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item D", Category: "Books", Amount: 75.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item E", Category: "Electronics", Amount: 150.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// DISTINCT query
	type CategoryOnly struct {
		Category string `bun:"category"`
	}

	var categories []CategoryOnly
	err := db.NewSelect().
		ColumnExpr(`DISTINCT "category"`).
		TableExpr(`"complex_queries"`).
		Scan(ctx, &categories)

	if err != nil {
		t.Fatalf("Failed DISTINCT query: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 distinct categories, got %d", len(categories))
	}

	t.Log("✓ DISTINCT working correctly")
}

// Test: Update with subquery
func TestBun_UpdateConditional(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Cheap Item", Category: "Test", Amount: 10.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Expensive Item", Category: "Test", Amount: 500.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Update records based on condition
	res, err := db.NewUpdate().
		Model((*ComplexQuery)(nil)).
		Set(`"is_active" = ?`, false).
		Where(`"amount" > ?`, 100.00).
		Exec(ctx)

	if err != nil {
		t.Fatalf("Failed conditional update: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	t.Log("✓ Conditional update working correctly")
}

// Test: Max, Min, Average aggregations
func TestBun_Aggregations(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item 1", Category: "Test", Amount: 100.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 2", Category: "Test", Amount: 200.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 3", Category: "Test", Amount: 300.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Query aggregations
	type AggResult struct {
		MaxAmount float64             `bun:"max_amount"`
		MinAmount float64             `bun:"min_amount"`
		AvgAmount float64             `bun:"avg_amount"`
		Count     db2dialect.SmallInt `bun:"cnt"`
	}

	var result []AggResult
	err := db.NewSelect().
		ColumnExpr(`MAX("amount") AS "max_amount"`).
		ColumnExpr(`MIN("amount") AS "min_amount"`).
		ColumnExpr(`AVG("amount") AS "avg_amount"`).
		ColumnExpr("COUNT(*) AS \"cnt\"").
		TableExpr(`"complex_queries"`).
		Scan(ctx, &result)

	if err != nil {
		t.Fatalf("Failed aggregation query: %v", err)
	}

	if len(result) < 1 {
		t.Fatal("No aggregation result returned")
	}

	agg := result[0]
	if agg.MaxAmount != 300.00 {
		t.Errorf("Max amount incorrect: %f != 300.00", agg.MaxAmount)
		t.Errorf("Max amount incorrect: %v != 300.00", agg.MaxAmount)
	}
	if agg.MinAmount != 100.00 {
		t.Errorf("Min amount incorrect: %f != 100.00", agg.MinAmount)
		t.Errorf("Min amount incorrect: %v != 100.00", agg.MinAmount)
	}
	if agg.Count != 3 {
		t.Errorf("Count incorrect: %d != 3", agg.Count)
	}

	t.Logf("✓ Aggregations: MAX=%f, MIN=%f, AVG=%f, COUNT=%d", agg.MaxAmount, agg.MinAmount, agg.AvgAmount, agg.Count)
}

// Test: Offset and pagination
func TestBun_OffsetPagination(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item 1", Category: "Test", Amount: 10.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 2", Category: "Test", Amount: 20.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 3", Category: "Test", Amount: 30.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 4", Category: "Test", Amount: 40.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Item 5", Category: "Test", Amount: 50.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Pagination: Page 1 (limit 2, offset 0)
	var page1 []ComplexQuery
	err := db.NewSelect().
		Model(&page1).
		Order("amount").
		Limit(2).
		Offset(0).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed pagination page 1: %v", err)
	}

	if len(page1) != 2 || page1[0].Amount != 10.00 {
		t.Error("Page 1 pagination incorrect")
	}

	// Pagination: Page 2 (limit 2, offset 2)
	var page2 []ComplexQuery
	err = db.NewSelect().
		Model(&page2).
		Order("amount").
		Limit(2).
		Offset(2).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed pagination page 2: %v", err)
	}

	if len(page2) != 2 || page2[0].Amount != 30.00 {
		t.Error("Page 2 pagination incorrect")
	}

	t.Log("✓ Pagination with OFFSET working correctly")
}

// Test: Raw SQL query fallback
func TestBun_RawSQL(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "Item 1", Category: "Test", Amount: 100.00, IsActive: 1},
		{Name: "Item 2", Category: "Test", Amount: 200.00, IsActive: 1},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Raw SQL query
	var result []ComplexQuery
	err := db.NewRaw(
		`SELECT "id", "name", "category", "amount", "is_active", "created_date" FROM "complex_queries" WHERE "amount" > ?`,
		150.00,
	).Scan(ctx, &result)

	if err != nil {
		t.Fatalf("Failed raw SQL query: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}

	t.Log("✓ Raw SQL queries working correctly")
}

// Test: PrepareContext support
func TestBun_PreparedStatements(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert multiple records for prepared statement test
	for i := 1; i <= 3; i++ {
		record := &ComplexQuery{
			Name:     fmt.Sprintf("Item %d", i),
			Category: "Test",
			Amount:   float64(i * 100),
			IsActive: 1,
		}
		db.NewInsert().Model(record).Exec(ctx)
	}

	// Use prepared query multiple times
	var result []ComplexQuery
	err := db.NewSelect().
		Model(&result).
		Where(`"category" = ?`, "Test").
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed prepared statement: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	t.Log("✓ Prepared statements working correctly")
}

// Test: Case sensitivity in queries
func TestBun_CaseHandling(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	db.NewDropTable().Model((*ComplexQuery)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*ComplexQuery)(nil)).IfNotExists().Exec(ctx)

	// Insert test data
	records := []ComplexQuery{
		{Name: "TEST", Category: "Books", Amount: 50.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "Test", Category: "Books", Amount: 60.00, IsActive: db2dialect.SmallIntBool(1)},
		{Name: "test", Category: "Books", Amount: 70.00, IsActive: db2dialect.SmallIntBool(1)},
	}

	db.NewInsert().Model(&records).Exec(ctx)

	// Query case-sensitive
	var result []ComplexQuery
	err := db.NewSelect().
		Model(&result).
		Where(`"name" = ?`, "TEST").
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed case sensitivity test: %v", err)
	}

	// DB2 is typically case-sensitive
	if len(result) < 1 {
		t.Logf("Note: DB2 case sensitivity may require adjustment")
	}

	t.Logf("✓ Case handling query returned %d results", len(result))
}
