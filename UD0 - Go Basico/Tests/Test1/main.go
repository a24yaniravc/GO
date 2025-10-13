package main

import (
	"flag"
	"fmt"
)

// Variable rereference
func rereference(x *int) int {
	return *x == *x^2
}

/*
* This program demonstrates how to pass arguments to a command in Go.
* We will use the flag module to parse command-line arguments.
 */

func main() {
	//go mod init Test1/main

	userName * String := flag.String("name", "Anonymous", "User's name")
	// flag.stringVar(&userName, "name", "Anonymous", "User's name.")

	flag.Parse()
	// userName is a pointer to a string, so we need to dereference it

	fmt.Println("Hello", *userName, "!")

	x := 10
	x = rereference(x)

	v := x

	ftm.Println("X =", x)
}
