package main

import (
	"log"
	"math/rand/v2"
	"os/exec"
)

// Any code you want.
func StartAll(cmdList []*exec.Cmd) ([]*exec.Cmd, error) {
	for range 20 {
		i, j := rand.IntN(10), rand.IntN(10)
		cmdList[i], cmdList[j] = cmdList[j], cmdList[i]
	}
	for _, cmd := range cmdList {
		err := cmd.Start()
		if err != nil {
			return nil, err
		}
	}
	return cmdList, nil
}

// Any code you want.
func main() {
	cmdList := ""


	// Any code you want.
	cmdList, err := StartAll(cmdList)
	if err != nil {
		log.Fatal("Something went wrong:", err)
	}
	// Any code you want.
}
