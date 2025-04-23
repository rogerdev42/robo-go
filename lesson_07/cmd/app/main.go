package main

import (
	"lesson_07/internal/documentstore"
)

func main() {
	defer documentstore.CloseLogFile()
}
