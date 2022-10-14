package main

import (
    "testing"
)
func TestConnectionPool(t *testing.T){
    if(ConnectionPool() == 0){
        t.Error("Error in Connection pool")
    }
}
