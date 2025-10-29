package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	// We will use the "flag" package to read command line arguments:
	ipAddr := flag.String("ip", "10.21.0.12", "IP address to listen on")
	port := flag.String("port", "8080", "Port number to listen on")

	// First we neet to listen for new connections:
	address := fmt.Sprintf("%s:%s", *ipAddr, *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// We need to remember to close the listener when we're done:
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

	// Now we can accept incoming connections in a loop:
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
		}

		// And read the connection for data:
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
		}

		fmt.Printf("Received a %d bytes message: %q\n", n, string(buffer[:n]))

		// We can send a response back to the client:
		n, err = conn.Write([]byte("Message received loud and clear!\n"))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
		}

		fmt.Printf("Sent a %d bytes response\n", n)

		// Now we can close the connection:
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %s\n", err)
		}

	}
}
