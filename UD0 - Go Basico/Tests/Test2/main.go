package main

import(
	"fmt"
	"bufio"
	"os"
)

// ReadLine
func main() {
	scanner := bufio.NewScanner(os.Stdin) 
	
	fmt.Print("Dime un número: ")
	scanner.Scan()

	fmt.Printf("El número escrito fue %s.\n", scanner.Text())
}