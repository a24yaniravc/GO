// primera

package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No numbers provided")
		os.Exit(1)
	}

	var nums []int
	for _, arg := range os.Args[1:] {
		n, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Error: '%s' is not a valid number\n", arg)
			os.Exit(1)
		}
		nums = append(nums, n)
	}

	// Ascending
	asc := append([]int{}, nums...)
	sort.Ints(asc)

	// Descending
	desc := append([]int{}, nums...)
	sort.Sort(sort.Reverse(sort.IntSlice(desc)))

	// Even-then-odd
	eo := append([]int{}, nums...)
	sort.Slice(eo, func(i, j int) bool {
		if eo[i]%2 == 0 && eo[j]%2 != 0 {
			return true
		}
		if eo[i]%2 != 0 && eo[j]%2 == 0 {
			return false
		}
		return eo[i] < eo[j]
	})

	fmt.Println("Ascending:")
	fmt.Println(sliceToString(asc))

	fmt.Println("Descending:")
	fmt.Println(sliceToString(desc))

	fmt.Println("By even and odd:")
	fmt.Println(sliceToString(eo))
}

func sliceToString(s []int) string {
	str := ""
	for i, v := range s {
		if i > 0 {
			str += " "
		}
		str += strconv.Itoa(v)
	}
	return str
}

// Segunda

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: capture <numbers...>")
		return
	}

	// Execute the orders program
	cmd := exec.Command("./orders", os.Args[1:]...)
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error running orders program:", err)
		os.Exit(1)
	}

	fmt.Println("Captured output:")
	fmt.Println(string(output))
}

// Tercera
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	// WRONG: Only one pipe used, insufficient for bidirectional communication
	// right := exec.Command("cat")

	// FIX: We need two separate pipes: parent→child and child→parent
	parentToChildR, parentToChildW, _ := os.Pipe()
	childToParentR, childToParentW, _ := os.Pipe()

	cmd := exec.Command(os.Args[0], "child")

	// WRONG: Child was inheriting stdin/stdout incorrectly
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout

	// FIX: Child receives one pipe as stdin and one pipe as stdout
	cmd.Stdin = parentToChildR
	cmd.Stdout = childToParentW

	// Start child process
	cmd.Start()

	message := "hello child process"
	fmt.Println("Parent: Sending message to child ->", message)

	// Send message
	parentToChildW.Write([]byte(message + "\n"))
	parentToChildW.Close()

	// Receive uppercase message
	buf := make([]byte, 256)
	n, _ := childToParentR.Read(buf)
	upper := strings.TrimSpace(string(buf[:n]))

	fmt.Println("Parent: Received from child ->", upper)

	// Send acknowledgement
	fmt.Println("Parent: Sending acknowledgement")
	parentAck, _ := os.Pipe()
	// WRONG: Child never received acknowledgement pipe
	// FIX: Use a simple method: send ack back through childToParent pipe (closing it)
	childToParentW.Close()

	cmd.Wait()
}

// ---------------- CHILD PROCESS ------------------

func child() {
	fmt.Println("Child: Started")

	// Read message from parent
	buf := make([]byte, 256)
	n, _ := os.Stdin.Read(buf)
	msg := strings.TrimSpace(string(buf[:n]))

	fmt.Println("Child: Received ->", msg)

	upper := strings.ToUpper(msg)

	fmt.Println("Child: Sending uppercase ->", upper)

	// Send uppercase string back
	fmt.Println(upper)

	// Wait for acknowledgement
	fmt.Println("Child: Waiting for acknowledgement...")
	_, err := os.Stdout.Write([]byte{})
	if err != nil {
		fmt.Println("Child: Acknowledged. Exiting.")
	}
}

func init() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		child()
		os.Exit(0)
	}
}

// Tercera (extra)
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Child mode
	if len(os.Args) > 1 && os.Args[1] == "child" {
		runChild()
		return
	}

	// -----------------------
	// PARENT PROCESS
	// -----------------------

	// WRONG: the original code tried to create pipes manually
	// and attach them incorrectly to stdin/stdout of the child.
	// (This section is intentionally left commented as required)
	/*
		pipe := os.Pipe()
		cmd := exec.Command("child")
		cmd.Stdin = os.Stdout  // wrong
		cmd.Stdout = pipe      // wrong
	*/

	// FIX: Use StdinPipe() and StdoutPipe(), the simplest correct approach
	cmd := exec.Command(os.Args[0], "child")

	childStdin, _ := cmd.StdinPipe()
	childStdout, _ := cmd.StdoutPipe()

	cmd.Start()

	message := "hello child"
	fmt.Println("Parent: Sending message to child ->", message)

	// Send message
	childStdin.Write([]byte(message + "\n"))

	// Read uppercase response
	reader := bufio.NewReader(childStdout)
	upper, _ := reader.ReadString('\n')
	upper = strings.TrimSpace(upper)

	fmt.Println("Parent: Received from child ->", upper)

	// Send ACK
	fmt.Println("Parent: Sending ACK")
	childStdin.Write([]byte("ACK\n"))
	childStdin.Close()

	cmd.Wait()
}

func runChild() {
	// -----------------------
	// CHILD PROCESS
	// -----------------------

	fmt.Println("Child: Started")

	reader := bufio.NewReader(os.Stdin)

	// Receive message
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	fmt.Println("Child: Received ->", msg)

	upper := strings.ToUpper(msg)

	// Send uppercase
	fmt.Println("Child: Sending uppercase ->", upper)
	fmt.Println(upper)

	// Wait for ACK
	fmt.Println("Child: Waiting for ACK...")
	ack, _ := reader.ReadString('\n')
	ack = strings.TrimSpace(ack)

	if ack == "ACK" {
		fmt.Println("Child: ACK received. Exiting.")
	} else {
		fmt.Println("Child: Invalid ACK. Exiting anyway.")
	}
}


// EXTRA
// Qué pasa realmente

// 1. Los argumentos de una función 'defer' se evalúan inmediatamente.
// captura x = 10 en ese momento.

// 2. Las variables usadas dentro de un closure defer NO se evalúan inmediatamente.
// Se evalúan cuando el defer se ejecuta.

// Entonces: 
// Defer 2 executed, x = 10
// Defer 1 executed, x = 20

// Defer 2 captura el valor 10 → porque los parámetros del defer se evalúan cuando se declara.
// Defer 1 ve el valor final de x, que es 20, porque usa la variable directamente.