package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

func panicHelper(err error) {
	fmt.Println(err)
	fmt.Println("stack trace for helping debug:")
	debug.PrintStack()
	os.Exit(1)
}
