package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"os/exec"
)

func main() {
	// Obtener la ruta del ejecutable del proceso hijo
	fmt.Println("Parent process started with PID:", os.Getpid())
	cmdPath := "/path/to/child/executable" // Replace with the actual path to the child process executable
	cmd := exec.Command(cmdPath)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return
	}
		defer stdinPipe.Close()

	// Crear un canal para recibir señales del sistema operativo
	go func ()  {
		canal := make(chan os.Signal, 1)
	signal.Notify(canal, syscall.SIGINT)

	switch <-canal {
	case syscall.SIGINT:
		fmt.Println("Parent process received SIGINT, sending SIGTERM to child process")
	case syscall.SIGTERM:
		fmt.Println("Parent process received SIGTERM, terminating child process")
	}
}()

	defer stdinPipe.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	
	defer stdoutPipe.Close()
}