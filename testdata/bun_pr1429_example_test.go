package main

// This test mirrors the example program submitted in
// https://github.com/uptrace/bun/pull/1429 (example/db2/main.go), adapted to
// use this repo's shared Createconnection()/GetDialect() test helpers instead
// of a hardcoded DSN.

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/uptrace/bun"
)

type PR1429User struct {
	bun.BaseModel `bun:"table:bun_db2_demo_users"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	Email     string    `bun:"email,notnull,unique"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func TestBun_PR1429_Example(t *testing.T) {
	ctx := context.Background()

	sqldb := Createconnection()
	if sqldb == nil {
		t.Fatal("Failed to create SQL connection")
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, GetDialect())
	defer db.Close()

	tableName := fmt.Sprintf("bun_db2_demo_users_%d", time.Now().Unix())
	table := bun.Ident(tableName)

	if _, err := db.NewCreateTable().Model((*PR1429User)(nil)).ModelTableExpr("?", table).Exec(ctx); err != nil {
		t.Fatalf("create table failed: %v", err)
	}
	defer func() {
		_, _ = db.NewDropTable().Model((*PR1429User)(nil)).ModelTableExpr("?", table).Exec(ctx)
	}()
	t.Logf("created table %s", strings.ToUpper(tableName))
	printPR1429Users(t, ctx, db, table, "after create")

	user := &PR1429User{
		Name:      "Alice",
		Email:     fmt.Sprintf("alice+%d@example.com", time.Now().UnixNano()),
		CreatedAt: time.Now(),
	}
	if _, err := db.NewInsert().Model(user).ModelTableExpr("?", table).Exec(ctx); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	t.Log("inserted demo user")
	printPR1429Users(t, ctx, db, table, "after insert")

	if _, err := db.NewDelete().Model((*PR1429User)(nil)).ModelTableExpr("?", table).Where("1 = 1").Exec(ctx); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	t.Log("deleted demo users")
	printPR1429Users(t, ctx, db, table, "after delete")

	if _, err := db.NewDropTable().Model((*PR1429User)(nil)).ModelTableExpr("?", table).Exec(ctx); err != nil {
		t.Fatalf("drop table failed: %v", err)
	}
	t.Logf("dropped table %s", strings.ToUpper(tableName))
}

func printPR1429Users(t *testing.T, ctx context.Context, db *bun.DB, table bun.Ident, label string) {
	var users []PR1429User
	if err := db.NewSelect().
		Model(&users).
		ModelTableExpr("? AS ?", table, bun.Ident("pr1429_user")).
		OrderExpr(`"pr1429_user"."id" ASC`).
		Limit(10).
		Scan(ctx); err != nil {
		t.Fatalf("select failed: %v", err)
	}
	t.Logf("%s: %v", label, users)
}
