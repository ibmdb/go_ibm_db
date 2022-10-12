package main

import (
    "testing"
)
func TestConnectionPoolWithTimeout(t *testing.T){
    if(ConnectionPoolWithTimeout() == 0){
        t.Error("Error in Connection pool with timeout")
    }
}
