package main

import (
    "testing"
)
func TestPoolOpen(t *testing.T){
    if(PoolOpen() == 0){
        t.Error("pool connection not opened")
    }
}