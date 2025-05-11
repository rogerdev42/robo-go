package main

import (
	"fmt"
	"lesson_11/internal/documentstore"
	"math/rand"
	"strconv"
	"sync"
)

func main() {
	defer documentstore.CloseLogFile()

	store := documentstore.NewStore()
	collection, _ := store.CreateCollection("test-store", &documentstore.CollectionConfig{
		PrimaryKey: "id",
	})
	err := collection.CreateIndex("id")
	if err != nil {
		fmt.Println("Error creating index:", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(goroutineId int) {
			defer wg.Done()

			documentId := strconv.Itoa(rand.Intn(100))

			switch rand.Intn(3) {
			case 0:
				err := collection.Put(documentstore.Document{
					Fields: map[string]documentstore.DocumentField{
						"id": {
							Type:  documentstore.DocumentFieldTypeString,
							Value: documentId,
						},
						"data": {
							Type:  documentstore.DocumentFieldTypeString,
							Value: fmt.Sprintf("goroutine-%d", goroutineId),
						},
					},
				})
				if err != nil {
					fmt.Println("Error creating document:", documentId, err)
				} else {
					fmt.Println("Document created successfully:", documentId)
				}
			case 1:
				_, err := collection.Query("id", documentstore.QueryParams{
					MinValue: &documentId,
					MaxValue: &documentId,
				})
				if err != nil {
					fmt.Println("Error querying documents:", documentId, err)
				} else {
					fmt.Println("Documents found successfully:", documentId)
				}
			case 2:
				err := collection.Delete(documentId)
				if err != nil {
					fmt.Println("Error deleting document:", documentId, err)
				} else {
					fmt.Println("Document deleted successfully:", documentId)
				}
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("That's all")
}
