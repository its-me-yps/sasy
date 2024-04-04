package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sayymeer/sasy/pkg/git"
)

func main() {
	args := os.Args[1:]

	// If no args are provided
	if len(args) == 0 {
		// TODO: To display help section when no args provided
		fmt.Printf("sasy: no args provided\n")
		os.Exit(0)
	}

	command := args[0]
	switch command {
	case "init":
		git.InitHandler()
	case "commit":
		git.CommitHandler()
	default:
		log.Printf("sasy: %s is not a valid command\n", command)
		os.Exit(1)
	}

}
