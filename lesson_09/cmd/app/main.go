package main

import (
	"lesson_09/internal/documentstore"
)

func main() {
	defer documentstore.CloseLogFile()
}
