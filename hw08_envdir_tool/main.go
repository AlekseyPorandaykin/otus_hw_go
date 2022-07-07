package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		return
	}
	envs, errRead := ReadDir(args[1])
	if errRead != nil {
		return
	}
	if code := RunCmd(args[2:], envs); code > 0 {
		fmt.Println("Error execute")
	}
}
