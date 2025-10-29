package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	// Lets check if we get an ip and port from command line:
	ipAddr := flag.String("ip", "10.21.0.12", "IP address to connect to")
	port := flag.String("port", "8080", "Port number to connect to")

	addr := fmt.Sprintf("%s:%s", *ipAddr, *port)

	// Now we'll try to establish a connection to the server:
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}

	// We can write something to the server:
	n, err := conn.Write([]byte("Hola, buenos dias :3c!\n"))
	if err != nil {
		log.Fatalf("Error writing to server: %v\n", err)
	}

	fmt.Printf("Sent a %d bytes message to server\n", n)

	// And check if the server has something to say back:
	buff := make([]byte, 1024)
	n, err = conn.Read(buff)
	if err != nil {
		log.Fatalf("Error reading from server: %v\n", err)
	}

	fmt.Printf("Received a %d bytes response from server: %q\n", n, string(buff[:n]))
}
