// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package odbc implements database/sql driver to access data via odbc interface.
package log2

import (
	"log"
	"os"
	"fmt"
)

var globalvar string = ""
var globalArgsLen int = 0

func GetPath(filename string, argsLen int) {
	//fmt.Println("filename = ", filename)
	//fmt.Println("Args Length = ", argsLen)

	if _, err := os.Stat(filename); err == nil {
        //fmt.Println("File exits\n")
		e := os.Remove(filename)
		if e != nil {
		  fmt.Println("Problem in removing existing log file")
		}
	}

	globalvar = filename
	globalArgsLen = argsLen

}

func Trace1(msg1 string) {
    if globalvar != "" {
		//file, errlog := os.OpenFile("C:\\temp\\testlogs2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		file, errlog := os.OpenFile(globalvar, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if errlog != nil {
			log.Fatal(errlog)
		}
		log.SetOutput(file)
		log.Println(msg1)
    } else if globalArgsLen > 1 {
	    log.SetOutput(os.Stdout)
		log.Println(msg1)
	}

}

