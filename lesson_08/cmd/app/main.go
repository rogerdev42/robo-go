package main

import (
	"lesson_08/internal/documentstore"
)

func main() {
	defer documentstore.CloseLogFile()
}
