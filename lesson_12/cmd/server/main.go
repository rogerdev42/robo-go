package main

import (
	"log"
)

func main() {
	if err := StartServer("localhost:8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
