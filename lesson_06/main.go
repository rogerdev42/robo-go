package main

import (
	"fmt"
	"lesson_06/documentstore"
)

func main() {
	defer documentstore.CloseLogFile()

	// 1. Create store and add test data
	store := fillStore()

	// 2. Get the dump of the store
	dump, err := store.Dump()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(dump)

	// 3. Put the given dump (prev. step result) into the getDump() result

	// 4. Create store from the dump
	store = documentstore.NewStore()
	dump = getDump()
	store, err = documentstore.NewStoreFromDump(dump)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(store.GetCollection("products"))

	// 4. Dump store to file

	err = store.DumpToFile("dump")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// 5. Create store from file
	store = documentstore.NewStore()
	store, err = documentstore.NewStoreFromFile("dump")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println(store.GetCollection("users"))
}

func fillStore() *documentstore.Store {
	store := documentstore.NewStore()
	cfg := documentstore.CollectionConfig{
		PrimaryKey: "ID",
	}

	userCollection, _ := store.CreateCollection("users", &cfg)

	user1, _ := documentstore.MarshalDocument(struct {
		ID    string
		Name  string
		Email string
	}{
		ID:    "1",
		Name:  "Alice",
		Email: "alice@example.com",
	})

	user2, _ := documentstore.MarshalDocument(struct {
		ID    string
		Name  string
		Email string
	}{
		ID:    "2",
		Name:  "Bob",
		Email: "bob@example.com",
	})

	userCollection.Put(*user1)
	userCollection.Put(*user2)

	productCollection, _ := store.CreateCollection("products", &cfg)

	product1, _ := documentstore.MarshalDocument(struct {
		ID    string
		Name  string
		Price float64
	}{
		ID:    "p1",
		Name:  "Laptop",
		Price: 1200.50,
	})

	product2, _ := documentstore.MarshalDocument(struct {
		ID    string
		Name  string
		Price float64
	}{
		ID:    "p2",
		Name:  "Smartphone",
		Price: 799.99,
	})

	productCollection.Put(*product1)
	productCollection.Put(*product2)

	return store
}

func getDump() []byte {
	return []byte{
		123, 34, 99, 111, 108, 108, 101, 99, 116, 105, 111, 110, 115, 34, 58, 123,
		34, 112, 114, 111, 100, 117, 99, 116, 115, 34, 58, 123, 34, 99, 111, 110,
		102, 105, 103, 34, 58, 123, 34, 80, 114, 105, 109, 97, 114, 121, 75, 101,
		121, 34, 58, 34, 73, 68, 34, 125, 44, 34, 100, 111, 99, 117, 109, 101,
		110, 116, 115, 34, 58, 123, 34, 112, 49, 34, 58, 123, 34, 70, 105, 101,
		108, 100, 115, 34, 58, 123, 34, 73, 68, 34, 58, 123, 34, 84, 121, 112,
		101, 34, 58, 34, 115, 116, 114, 105, 110, 103, 34, 44, 34, 86, 97, 108,
		117, 101, 34, 58, 34, 112, 49, 34, 125, 44, 34, 78, 97, 109, 101, 34,
		58, 123, 34, 84, 121, 112, 101, 34, 58, 34, 115, 116, 114, 105, 110, 103,
		34, 44, 34, 86, 97, 108, 117, 101, 34, 58, 34, 76, 97, 112, 116, 111,
		112, 34, 125, 44, 34, 80, 114, 105, 99, 101, 34, 58, 123, 34, 84, 121,
		112, 101, 34, 58, 34, 110, 117, 109, 98, 101, 114, 34, 44, 34, 86, 97,
		108, 117, 101, 34, 58, 49, 50, 48, 48, 46, 53, 125, 125, 125, 44, 34,
		112, 50, 34, 58, 123, 34, 70, 105, 101, 108, 100, 115, 34, 58, 123, 34,
		73, 68, 34, 58, 123, 34, 84, 121, 112, 101, 34, 58, 34, 115, 116, 114,
		105, 110, 103, 34, 44, 34, 86, 97, 108, 117, 101, 34, 58, 34, 112, 50,
		34, 125, 44, 34, 78, 97, 109, 101, 34, 58, 123, 34, 84, 121, 112, 101,
		34, 58, 34, 115, 116, 114, 105, 110, 103, 34, 44, 34, 86, 97, 108, 117,
		101, 34, 58, 34, 83, 109, 97, 114, 116, 112, 104, 111, 110, 101, 34, 125,
		44, 34, 80, 114, 105, 99, 101, 34, 58, 123, 34, 84, 121, 112, 101, 34,
		58, 34, 110, 117, 109, 98, 101, 114, 34, 44, 34, 86, 97, 108, 117, 101,
		34, 58, 55, 57, 57, 46, 57, 57, 125, 125, 125, 125, 125, 44, 34, 117,
		115, 101, 114, 115, 34, 58, 123, 34, 99, 111, 110, 102, 105, 103, 34, 58,
		123, 34, 80, 114, 105, 109, 97, 114, 121, 75, 101, 121, 34, 58, 34, 73,
		68, 34, 125, 44, 34, 100, 111, 99, 117, 109, 101, 110, 116, 115, 34, 58,
		123, 34, 49, 34, 58, 123, 34, 70, 105, 101, 108, 100, 115, 34, 58, 123,
		34, 69, 109, 97, 105, 108, 34, 58, 123, 34, 84, 121, 112, 101, 34, 58,
		34, 115, 116, 114, 105, 110, 103, 34, 44, 34, 86, 97, 108, 117, 101, 34,
		58, 34, 97, 108, 105, 99, 101, 64, 101, 120, 97, 109, 112, 108, 101, 46,
		99, 111, 109, 34, 125, 44, 34, 73, 68, 34, 58, 123, 34, 84, 121, 112,
		101, 34, 58, 34, 115, 116, 114, 105, 110, 103, 34, 44, 34, 86, 97, 108,
		117, 101, 34, 58, 34, 49, 34, 125, 44, 34, 78, 97, 109, 101, 34, 58,
		123, 34, 84, 121, 112, 101, 34, 58, 34, 115, 116, 114, 105, 110, 103, 34,
		44, 34, 86, 97, 108, 117, 101, 34, 58, 34, 65, 108, 105, 99, 101, 34,
		125, 125, 125, 44, 34, 50, 34, 58, 123, 34, 70, 105, 101, 108, 100, 115,
		34, 58, 123, 34, 69, 109, 97, 105, 108, 34, 58, 123, 34, 84, 121, 112,
		101, 34, 58, 34, 115, 116, 114, 105, 110, 103, 34, 44, 34, 86, 97, 108,
		117, 101, 34, 58, 34, 98, 111, 98, 64, 101, 120, 97, 109, 112, 108, 101,
		46, 99, 111, 109, 34, 125, 44, 34, 73, 68, 34, 58, 123, 34, 84, 121,
		112, 101, 34, 58, 34, 115, 116, 114, 105, 110, 103, 34, 44, 34, 86, 97,
		108, 117, 101, 34, 58, 34, 50, 34, 125, 44, 34, 78, 97, 109, 101, 34,
		58, 123, 34, 84, 121, 112, 101, 34, 58, 34, 115, 116, 114, 105, 110, 103,
		34, 44, 34, 86, 97, 108, 117, 101, 34, 58, 34, 66, 111, 98, 34, 125,
		125, 125, 125, 125, 125, 125}
}
