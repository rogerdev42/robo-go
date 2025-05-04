package main

import (
	"fmt"
	"lesson_08/internal/documentstore"
	"lesson_08/lru"
)

func main() {
	defer documentstore.CloseLogFile()

	cache := lru.NewLruCache(2)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Put("key1", "value11")
	cache.Put("key3", "value3")

	val, ok := cache.Get("key3")
	fmt.Println(val, ok)

}
