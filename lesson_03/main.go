package main

import (
	"fmt"
	"lesson_3/document_store"
)

func main() {
	/*
		1. Positive scenario
	*/
	document := document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"key": {
				Type:  document_store.DocumentFieldTypeString,
				Value: "doc1",
			},
			"title": {
				Type:  document_store.DocumentFieldTypeString,
				Value: "Sample Document",
			},
			"views": {
				Type:  document_store.DocumentFieldTypeNumber,
				Value: 123,
			},
			"deleted": {
				Type:  document_store.DocumentFieldTypeBool,
				Value: false,
			},
			"array": {
				Type:  document_store.DocumentFieldTypeArray,
				Value: make([]int, 0),
			},
			"object": {
				Type:  document_store.DocumentFieldTypeObject,
				Value: make(map[string]interface{}),
			},
		},
	}

	// Get gocument list (empty)
	documentList := document_store.List()
	fmt.Println(documentList)

	// Put the document into store
	if ok := document_store.Put(document); ok {
		fmt.Println("Document added")
	} else {
		fmt.Println("Document not added")
	}

	// Get gocument list again
	documentList = document_store.List()
	fmt.Println(documentList)

	if doc, ok := document_store.Get("doc1"); ok {
		fmt.Println(*doc)
	} else {
		fmt.Println("Document not found")
	}

	// Delete the document
	if deleted := document_store.Delete("doc1"); deleted {
		fmt.Println("Document deleted")
	} else {
		fmt.Println("Document not found")
	}

	// Get gocument list (empty)
	documentList = document_store.List()

	/*
		2. Bad document (missing "key")
	*/
	document = document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"title": {
				Type:  document_store.DocumentFieldTypeString,
				Value: "Sample Document",
			},
			"views": {
				Type:  document_store.DocumentFieldTypeNumber,
				Value: 123,
			},
			"deleted": {
				Type:  document_store.DocumentFieldTypeBool,
				Value: false,
			},
			"array": {
				Type:  document_store.DocumentFieldTypeArray,
				Value: make([]int, 0),
			},
			"object": {
				Type:  document_store.DocumentFieldTypeObject,
				Value: make(map[string]interface{}),
			},
		},
	}

	// Put the document into store
	if ok := document_store.Put(document); ok {
		fmt.Println("Document added")
	} else {
		fmt.Println("Document not added")
	}

	/*
		3. Bad document (incorrect "key")
	*/
	document = document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"key": {
				Type:  document_store.DocumentFieldTypeNumber,
				Value: 123,
			},
			"title": {
				Type:  document_store.DocumentFieldTypeString,
				Value: "Sample Document",
			},
			"views": {
				Type:  document_store.DocumentFieldTypeNumber,
				Value: 123,
			},
			"deleted": {
				Type:  document_store.DocumentFieldTypeBool,
				Value: false,
			},
			"array": {
				Type:  document_store.DocumentFieldTypeArray,
				Value: make([]int, 0),
			},
			"object": {
				Type:  document_store.DocumentFieldTypeObject,
				Value: make(map[string]interface{}),
			},
		},
	}

	// Put the document into store
	if ok := document_store.Put(document); ok {
		fmt.Println("Document added")
	} else {
		fmt.Println("Document not added")
	}

	/*
		3. Get non-existent document
	*/
	if doc, ok := document_store.Get("xyz"); ok {
		fmt.Println(*doc)
	} else {
		fmt.Println("Document not found")
	}

	/*
		3. Delete non-existent document
	*/
	// Delete the document
	if deleted := document_store.Delete("xyz"); deleted {
		fmt.Println("Document deleted")
	} else {
		fmt.Println("Document not found")
	}

}
