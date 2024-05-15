package main

import (
	"fmt"
	"os"

	"sasy/pkg/sasy"
)

func main() {
	args := os.Args[1:]

	// If no args are provided
	if len(args) == 0 {
		// TODO: To display help section when no args provided
		fmt.Println("sasy: no commands provided")
		fmt.Println(sasy.Usage())
		os.Exit(1)
	}

	cmd, ok := sasy.Commands[args[0]]
	if !ok {
		fmt.Println("sasy: not a valid command")
		fmt.Println(sasy.Usage())
		os.Exit(1)
	}

	if err := cmd(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
