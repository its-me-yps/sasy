package main

import (
	"log"
	"os"

	"github.com/sayymeer/sasy/pkg/git"
)

func main() {
	args := os.Args[1:]
	command := args[0]
	switch command {
	case "init":
		git.InitHandler()
	case "commit":
		git.CommitHandler()
	default:
		log.Printf("sasy: %s is not a valid command", command)
		os.Exit(1)
	}

}
