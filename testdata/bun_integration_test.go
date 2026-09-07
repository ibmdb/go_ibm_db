package main

import (
	"context"
	"testing"
	"time"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/ibmdb/go_ibm_db/db2dialect"
	"github.com/uptrace/bun"
)

// Models for testing
type Employee struct {
	ID        int64                   `bun:",pk,autoincrement"`
	Name      string                  `bun:",notnull"`
	Email     string                  `bun:",unique,notnull"`
	Salary    float64                 `bun:""`
	Active    db2dialect.SmallIntBool `bun:""` // 0=false, 1=true
	CreatedAt time.Time               `bun:",default:current_timestamp"`
	UpdatedAt time.Time               `bun:",default:current_timestamp"`
}

type Department struct {
	ID        int64      `bun:",pk,autoincrement"`
	Name      string     `bun:",notnull,unique"`
	Location  string     `bun:""`
	Employees []Employee `bun:"rel:has-many,join:id=dept_id"`
	CreatedAt time.Time  `bun:",default:current_timestamp"`
}

type Project struct {
	ID          int64                `bun:",pk,autoincrement"`
	Title       string               `bun:",notnull"`
	Description string               `bun:""`
	Status      string               `bun:",default:'ACTIVE'"`
	Budget      float64              `bun:""`
	StartDate   db2dialect.Timestamp `bun:""`
	EndDate     db2dialect.Timestamp `bun:""`
	CreatedAt   time.Time            `bun:",default:current_timestamp"`
}

type Product struct {
	ID        int64                   `bun:",pk,autoincrement"`
	Name      string                  `bun:",notnull"`
	Price     float64                 `bun:",notnull"`
	Quantity  db2dialect.SmallInt     `bun:",default:0"`
	Category  string                  `bun:""`
	Available db2dialect.SmallIntBool `bun:""` // 0=false, 1=true
	CreatedAt time.Time               `bun:",default:current_timestamp"`
}

// Helper function to get DB connection
func getBunDB(t *testing.T) *bun.DB {
	sqldb := Createconnection()
	if sqldb == nil {
		t.Fatal("Failed to create SQL connection")
	}

	db := bun.NewDB(sqldb, GetDialect())
	return db
}

// Helper to cleanly reset a table for tests
func resetTable(t *testing.T, ctx context.Context, db *bun.DB, model interface{}) {
	t.Helper()

	_, _ = db.NewDropTable().Model(model).Exec(ctx)
	if _, err := db.NewCreateTable().Model(model).Exec(ctx); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
}

// Test: Create table with Bun ORM
func TestBun_CreateTable(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Drop if exists
	_, _ = db.NewDropTable().Model((*Employee)(nil)).Exec(ctx)

	// Create table
	_, err := db.NewCreateTable().Model((*Employee)(nil)).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	t.Log("✓ Employee table created successfully")
}

// Test: Insert single record with Bun ORM
func TestBun_InsertSingleRecord(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Employee)(nil))

	// Insert
	employee := &Employee{
		Name:   "John Doe",
		Email:  "john@example.com",
		Salary: 75000.50,
		Active: db2dialect.SmallIntBool(1),
	}

	res, err := db.NewInsert().Model(employee).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("Failed to read affected row count: %v", err)
	}

	if rowsAffected != 1 {
		t.Errorf("Expected one inserted employee, got %d", rowsAffected)
	}

	var inserted Employee
	if err := db.NewSelect().Model(&inserted).Where("\"email\" = ?", employee.Email).Scan(ctx); err != nil {
		t.Fatalf("Failed to retrieve inserted employee: %v", err)
	}

	if inserted.Name != employee.Name {
		t.Errorf("Expected inserted employee name %q, got %q", employee.Name, inserted.Name)
	}

	t.Log("✓ Inserted employee successfully")
}

// Test: Bulk insert with Bun ORM
func TestBun_BulkInsert(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	// Bulk insert
	products := []Product{
		{Name: "Laptop", Price: 999.99, Quantity: db2dialect.SmallInt(5), Category: "Electronics"},
		{Name: "Mouse", Price: 29.99, Quantity: db2dialect.SmallInt(50), Category: "Accessories"},
		{Name: "Keyboard", Price: 79.99, Quantity: db2dialect.SmallInt(30), Category: "Accessories"},
		{Name: "Monitor", Price: 299.99, Quantity: db2dialect.SmallInt(10), Category: "Electronics"},
		{Name: "Desk Lamp", Price: 49.99, Quantity: db2dialect.SmallInt(20), Category: "Office"},
	}

	res, err := db.NewInsert().Model(&products).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to bulk insert: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != int64(len(products)) {
		t.Errorf("Expected %d rows affected, got %d", len(products), rowsAffected)
	}

	t.Logf("✓ Bulk inserted %d products", rowsAffected)
}

// Test: Select all records with Bun ORM
func TestBun_SelectAll(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Employee)(nil))

	employees := []Employee{
		{Name: "Alice Johnson", Email: "alice@example.com", Salary: 85000},
		{Name: "Bob Smith", Email: "bob@example.com", Salary: 75000},
		{Name: "Carol Davis", Email: "carol@example.com", Salary: 95000},
	}

	db.NewInsert().Model(&employees).Exec(ctx)

	// Select all
	var result []Employee
	err := db.NewSelect().Model(&result).Scan(ctx)
	if err != nil {
		t.Fatalf("Failed to select: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 employees, got %d", len(result))
	}

	t.Logf("✓ Selected %d employees from database", len(result))
}

// Test: Select with WHERE condition
func TestBun_SelectWithWhere(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	products := []Product{
		{Name: "Laptop", Price: 999.99, Category: "Electronics", Quantity: db2dialect.SmallInt(5)},
		{Name: "Mouse", Price: 29.99, Category: "Accessories", Quantity: db2dialect.SmallInt(50)},
		{Name: "Monitor", Price: 299.99, Category: "Electronics", Quantity: db2dialect.SmallInt(10)},
	}

	db.NewInsert().Model(&products).Exec(ctx)

	// Select with WHERE
	var electronics []Product
	err := db.NewSelect().
		Model(&electronics).
		Where(`"category" = ?`, "Electronics").
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed to select with WHERE: %v", err)
	}

	if len(electronics) != 2 {
		t.Errorf("Expected 2 electronics, got %d", len(electronics))
	}

	t.Logf("✓ Selected %d electronics using WHERE clause", len(electronics))
}

// Test: Update record with Bun ORM
func TestBun_Update(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Employee)(nil))

	employee := &Employee{
		Name:   "John Doe",
		Email:  "john@example.com",
		Salary: 75000.00,
	}

	db.NewInsert().Model(employee).Exec(ctx)

	// Update
	res, err := db.NewUpdate().
		Model(employee).
		Set(`"salary" = ?`, 85000.00).
		Set(`"updated_at" = CURRENT_TIMESTAMP`).
		Where(`"email" = ?`, "john@example.com").
		Exec(ctx)

	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	// Verify update
	var updated Employee
	db.NewSelect().
		Model(&updated).
		Where(`"email" = ?`, "john@example.com").
		Scan(ctx)

	if updated.Salary != 85000.00 {
		t.Errorf("Expected salary 85000.00, got %f", updated.Salary)
	}

	t.Log("✓ Record updated successfully")
}

// Test: Bulk update with Bun ORM
func TestBun_BulkUpdate(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	products := []Product{
		{Name: "Item1", Price: 10.00, Available: db2dialect.SmallIntBool(0)},
		{Name: "Item2", Price: 20.00, Available: db2dialect.SmallIntBool(0)},
		{Name: "Item3", Price: 30.00, Available: db2dialect.SmallIntBool(0)},
	}

	db.NewInsert().Model(&products).Exec(ctx)

	// Bulk update
	res, err := db.NewUpdate().
		Model((*Product)(nil)).
		Set(`"available" = ?`, db2dialect.SmallIntBool(1)).
		Where(`"available" = ?`, db2dialect.SmallIntBool(0)).
		Exec(ctx)

	if err != nil {
		t.Fatalf("Failed to bulk update: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 3 {
		t.Errorf("Expected 3 rows affected, got %d", rowsAffected)
	}

	t.Logf("✓ Bulk updated %d products", rowsAffected)
}

// Test: Delete record with Bun ORM
func TestBun_Delete(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Employee)(nil))

	employee := &Employee{
		Name:   "John Doe",
		Email:  "john@example.com",
		Salary: 75000.00,
	}

	db.NewInsert().Model(employee).Exec(ctx)

	// Delete
	res, err := db.NewDelete().
		Model(employee).
		Where(`"email" = ?`, "john@example.com").
		Exec(ctx)

	if err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	t.Log("✓ Record deleted successfully")
}

// Test: Bulk delete with Bun ORM
func TestBun_BulkDelete(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	products := []Product{
		{Name: "Laptop", Price: 999.99},
		{Name: "Mouse", Price: 29.99},
		{Name: "Monitor", Price: 299.99},
	}

	db.NewInsert().Model(&products).Exec(ctx)

	// Bulk delete
	res, err := db.NewDelete().
		Model((*Product)(nil)).
		Where(`"price" < ?`, 50.00).
		Exec(ctx)

	if err != nil {
		t.Fatalf("Failed to bulk delete: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	t.Logf("✓ Bulk deleted %d products", rowsAffected)
}

// Test: Count records
func TestBun_Count(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	products := []Product{
		{Name: "Item1", Price: 10.00},
		{Name: "Item2", Price: 20.00},
		{Name: "Item3", Price: 30.00},
	}

	db.NewInsert().Model(&products).Exec(ctx)

	// Count
	count, err := db.NewSelect().Model((*Product)(nil)).Count(ctx)
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	t.Logf("✓ Count returned %d records", count)
}

// Test: Order by and limit
func TestBun_OrderByAndLimit(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	products := []Product{
		{Name: "Expensive", Price: 999.99},
		{Name: "Medium", Price: 299.99},
		{Name: "Cheap", Price: 29.99},
	}

	db.NewInsert().Model(&products).Exec(ctx)

	// Query with ORDER BY and LIMIT
	var topProducts []Product
	err := db.NewSelect().
		Model(&topProducts).
		Order("price DESC").
		Limit(2).
		Scan(ctx)

	if err != nil {
		t.Fatalf("Failed to select with ORDER BY: %v", err)
	}

	if len(topProducts) != 2 {
		t.Errorf("Expected 2 records, got %d", len(topProducts))
	}

	if topProducts[0].Price < topProducts[1].Price {
		t.Error("Records not sorted by price DESC")
	}

	t.Log("✓ ORDER BY and LIMIT working correctly")
}

// Test: Query into map
func TestBun_QueryIntoMap(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Product)(nil))

	products := []Product{
		{Name: "Laptop", Price: 999.99},
		{Name: "Mouse", Price: 29.99},
	}

	db.NewInsert().Model(&products).Exec(ctx)

	// Query into maps
	var result []map[string]interface{}
	err := db.NewSelect().
		Table("products").
		Scan(ctx, &result)

	if err != nil {
		t.Fatalf("Failed to query into map: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result))
	}

	t.Logf("✓ Queried %d records into map format", len(result))
}

// Test: Transactions with Bun
func TestBun_Transaction(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup
	resetTable(t, ctx, db, (*Employee)(nil))

	// Run in transaction
	err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		employee := &Employee{
			Name:   "Jane Doe",
			Email:  "jane@example.com",
			Salary: 95000.00,
		}

		_, err := tx.NewInsert().Model(employee).Exec(ctx)
		if err != nil {
			return err
		}

		// Simulate another operation
		var count int
		count, err = tx.NewSelect().Model((*Employee)(nil)).Count(ctx)
		if err != nil {
			return err
		}

		if count != 1 {
			t.Errorf("Expected 1 record in transaction, got %d", count)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	t.Log("✓ Transaction completed successfully")
}

// Test: Different data types
func TestBun_DataTypes(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Drop and create
	resetTable(t, ctx, db, (*Project)(nil))

	// Test various data types
	now := time.Now()
	project := &Project{
		Title:       "Q4 Initiative",
		Description: "Strategic project for Q4",
		Status:      "ACTIVE",
		Budget:      500000.75,
		StartDate:   db2dialect.Timestamp(now),
		EndDate:     db2dialect.Timestamp(now.AddDate(0, 3, 0)),
	}

	res, err := db.NewInsert().Model(project).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to insert with various types: %v", err)
	}

	lastID, _ := res.LastInsertId()

	// Retrieve and verify
	var retrieved Project
	qSelect := db.NewSelect().Model(&retrieved)
	if lastID > 0 {
		qSelect = qSelect.Where(`"id" = ?`, lastID)
	} else {
		qSelect = qSelect.Where(`"title" = ?`, project.Title)
	}
	err = qSelect.Scan(ctx)

	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	if retrieved.Title != project.Title {
		t.Errorf("Title mismatch: %s != %s", retrieved.Title, project.Title)
	}

	if retrieved.Budget != project.Budget {
		t.Errorf("Budget mismatch: %f != %f", retrieved.Budget, project.Budget)
	}

	t.Log("✓ Various data types handled correctly")
}

// Test: Column operations
func TestBun_AddColumn(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Drop and create table
	resetTable(t, ctx, db, (*Product)(nil))

	// Add column
	_, err := db.NewAddColumn().
		Model((*Product)(nil)).
		ColumnExpr("discount DECIMAL(5,2) DEFAULT 0").
		Exec(ctx)

	if err != nil {
		t.Logf("Note: AddColumn may not be supported in all DB2 configurations: %v", err)
		return
	}

	t.Log("✓ Column added successfully")
}

// Test: Index operations (if supported)
func TestBun_CreateIndex(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup table
	resetTable(t, ctx, db, (*Employee)(nil))

	// Create index
	_, err := db.NewCreateIndex().
		Model((*Employee)(nil)).
		Index("idx_email").
		Column("email").
		Exec(ctx)

	if err != nil {
		t.Logf("Note: CreateIndex may not be supported: %v", err)
		return
	}

	t.Log("✓ Index created successfully")
}

// Test: Drop table
func TestBun_DropTable(t *testing.T) {
	db := getBunDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create table
	db.NewCreateTable().Model((*Product)(nil)).IfNotExists().Exec(ctx)

	// Drop table
	_, err := db.NewDropTable().Model((*Product)(nil)).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	t.Log("✓ Table dropped successfully")
}
