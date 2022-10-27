package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatal("No port number was written!")
	}
	if args[1] == "1" {
		CreateServer(args)
	} else {
		ConnectToServer(args)
	}
}
