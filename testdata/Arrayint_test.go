package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestIntArray(t *testing.T) {
	if IntArray() != nil {
		t.Error("Error at IntArray")
	}
}

// IntArray function performs inserting int,int8,int16,int32,int64 datatypes.
func IntArray() error {
	db := Createconnection()
	defer db.Close()

	db.Exec("Drop table arr")
	_, err := db.Exec("create table arr(var1 int)")
	if err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}

	a := []int{2, 3}
	b := []int8{2, 3}
	c := []int16{2, 3}
	d := []int32{2, 3}
	e := []int64{2, 3}
	st, err := db.Prepare("Insert into arr values(?)")
	if err != nil {
		fmt.Println("Prepare error: ", err)
		return err
	}
	defer st.Close()
	_, err = st.Query(a)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int")
		return err
	}
	_, err = st.Query(b)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int8")
		return err
	}
	_, err = st.Query(c)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int16")
		return err
	}
	_, err = st.Query(d)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int32")
		return err
	}
	_, err = st.Query(e)
	if !strings.Contains(fmt.Sprint(err), "did not create a result set") {
		fmt.Println("Error while inserting []int64")
		return err
	}
	return nil
}

func TestIntArrayThenScalar(t *testing.T) {
	db := Createconnection()
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	table := fmt.Sprintf("GIBM_PS_%d", time.Now().UnixNano())
	if _, err := conn.ExecContext(ctx, "CREATE TABLE "+table+" (id INTEGER)"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := conn.ExecContext(context.Background(), "DROP TABLE "+table); err != nil {
			t.Errorf("drop table %s: %v", table, err)
		}
	}()

	stmt, err := conn.PrepareContext(ctx, "INSERT INTO "+table+" (id) VALUES (?)")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	values := []interface{}{
		int64(0),
		[]int64{1, 2},
		int64(3),
		[]int64{4, 5, 6},
		int64(7),
		[]int64{8},
		int64(9),
	}
	for _, value := range values {
		if _, err := stmt.ExecContext(ctx, value); err != nil {
			t.Fatalf("insert %v: %v", value, err)
		}
	}

	var count, sum int
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*), SUM(id) FROM "+table).Scan(&count, &sum); err != nil {
		t.Fatal(err)
	}
	if count != 10 || sum != 45 {
		t.Fatalf("got count=%d and sum=%d, want count=10 and sum=45", count, sum)
	}
}
