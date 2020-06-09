package src

import (
	"fmt"
	"os"
)

/*
simple debug wrapper for debugging interfaces
*/
func Log(msg ...interface{}) {
	fmt.Println(msg...)
}

/*
if the assertion condition is false, then we have an error
*/
func Assert(assertion bool, msg string) {
	if !assertion {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
}
