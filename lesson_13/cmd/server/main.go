package main

import (
	"log"
)

func main() {
	if err := StartServer(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
