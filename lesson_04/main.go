package main

import (
	"fmt"
	"lesson_04/documentstore"
)

func main() {
	store := documentstore.NewStore()

	cfg := documentstore.CollectionConfig{
		PrimaryKey: "key",
	}

	ok, collection := store.CreateCollection("first", &cfg)
	if !ok {
		return
	}

	document := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "doc1",
			},
			"title": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Sample Document",
			},
			"views": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 123,
			},
			"deleted": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: false,
			},
			"array": {
				Type:  documentstore.DocumentFieldTypeArray,
				Value: make([]int, 0),
			},
			"object": {
				Type:  documentstore.DocumentFieldTypeObject,
				Value: make(map[string]interface{}),
			},
		},
	}

	collection.Put(document)

	fmt.Println(collection.List())

	col, ok := store.GetCollection("first")
	if !ok {
		return
	}
	fmt.Println(col)

	doc, ok := collection.Get("doc1")
	if !ok {
		return
	}

	fmt.Println(doc)

	collection.Delete("doc1")

	store.DeleteCollection("first")

	fmt.Println(collection.List())
}
