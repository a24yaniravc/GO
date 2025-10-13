package main

import (
	"fmt"
	"os/exec"
)

func main() {
	lsCmd := exec.Command("powershell.exe", "-c", "Get-ChildItem")
	fmt.Println("lsCmd:", lsCmd)

	
}