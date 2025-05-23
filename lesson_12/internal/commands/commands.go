package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"lesson_12/internal/documentstore"
	"log"
)

const (
	StoreStorageName     string = "store"
	CreateColCommandName string = "create" // Create a collection in the store
	ListColCommandName   string = "list"   // List documents in the collection
	DeleteColCommandName string = "delete" // Delete a collection from the store

	ColStorageName       string = "collection"
	PutDocCommandName    string = "put"    // Put a document in the collection
	GetDocCommandName    string = "get"    // Get a document from the collection
	DeleteDocCommandName string = "delete" // Delete a document from the collection
)

var store = documentstore.NewStore()

type CreateColCommandRequestPayload struct {
	Name   string                          `json:"name"`   // Collection Name
	Config *documentstore.CollectionConfig `json:"config"` // Collection Config
}
type DeleteColCommandRequestPayload struct {
	Name string `json:"name"` // Collection Name
}

type ListColCommandRequestPayload struct {
	Name string `json:"name"` // Collection Name
}
type PutDocCommandNameRequestPayload struct {
	Name string                 `json:"name"` // Collection Name
	Doc  map[string]interface{} `json:"doc"`
}

type GetDocCommandNameRequestPayload struct {
	Name string `json:"name"` // Collection Name
	Id   string `json:"Id"`   // Document ID
}

type DeleteDocCommandNameRequestPayload struct {
	Name string `json:"name"` // Collection Name
	Id   string `json:"Id"`   // Document ID
}

type CreateColCommandResponsePayload struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

type DeleteColCommandResponsePayload struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

type ListColCommandResponsePayload struct {
	Status string                   `json:"status"`
	Value  string                   `json:"value"`
	Docs   []documentstore.Document `json:"docs"`
}

type PutDocCommandNameResponsePayload struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

type GetDocCommandNameResponsePayload struct {
	Status string                 `json:"status"`
	Value  string                 `json:"value"`
	Doc    documentstore.Document `json:"doc"`
}

type DeleteDocCommandNameResponsePayload struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

func ExecStore(command string, param string) (resp string, err error) {
	switch command {
	case CreateColCommandName:
		resp, err = ExecCreateCol(param)
	case DeleteColCommandName:
		resp, err = ExecDeleteCol(param)
	case ListColCommandName:
		resp, err = ExecListCol(param)
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
	return
}

func ExecCol(command string, param string) (resp string, err error) {
	switch command {
	case PutDocCommandName:
		resp, err = ExecPutDoc(param)
	case GetDocCommandName:
		resp, err = ExecGetDoc(param)
	case DeleteDocCommandName:
		resp, err = ExecDeleteDoc(param)
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
	return
}

func ExecCreateCol(param string) (string, error) {
	p := &CreateColCommandRequestPayload{}
	err := json.Unmarshal([]byte(param), p)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}
	_, err = store.CreateCollection(p.Name, &documentstore.CollectionConfig{PrimaryKey: p.Config.PrimaryKey})
	if err != nil {
		return "", fmt.Errorf("collection creation error: %w", err)
	}

	r := CreateColCommandResponsePayload{
		Status: "created",
		Value:  p.Name,
	}
	resp, merr := json.Marshal(r)
	if merr != nil {
		log.Println("internal error:", merr)
		return "", errors.New("internal error")
	}
	return string(resp), nil
}

func ExecListCol(param string) (string, error) {
	p := &ListColCommandRequestPayload{}
	err := json.Unmarshal([]byte(param), p)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	col, err := store.GetCollection(p.Name)
	if err != nil {
		return "", fmt.Errorf("collection getting error: %w", err)
	}

	r := ListColCommandResponsePayload{
		Status: "success",
		Value:  p.Name,
		Docs:   col.List(),
	}

	resp, merr := json.Marshal(r)
	if merr != nil {
		log.Println("internal error:", merr)
		return "", errors.New("internal error")
	}
	return string(resp), nil
}

func ExecDeleteCol(param string) (string, error) {
	p := &DeleteColCommandRequestPayload{}
	err := json.Unmarshal([]byte(param), p)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}
	err = store.DeleteCollection(p.Name)
	if err != nil {
		return "", fmt.Errorf("collection deleting error: %w", err)
	}

	r := DeleteColCommandResponsePayload{
		Status: "deleted",
		Value:  p.Name,
	}
	resp, merr := json.Marshal(r)
	if merr != nil {
		log.Println("internal error:", merr)
		return "", errors.New("internal error")
	}
	return string(resp), nil
}

func ExecPutDoc(param string) (string, error) {
	p := &PutDocCommandNameRequestPayload{}
	err := json.Unmarshal([]byte(param), p)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	doc, merr := documentstore.MarshalDocument(p.Doc)
	if merr != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	collection, cerr := store.GetCollection(p.Name)
	if cerr != nil {
		return "", fmt.Errorf("collection getting error: %w", err)
	}

	err = collection.Put(*doc)
	if err != nil {
		return "", fmt.Errorf("put error: %w", err)
	}

	r := PutDocCommandNameResponsePayload{
		Status: "success",
	}

	resp, merr := json.Marshal(r)
	if merr != nil {
		log.Println("internal error:", merr)
		return "", errors.New("internal error")
	}
	return string(resp), nil

}

func ExecGetDoc(param string) (string, error) {
	p := &GetDocCommandNameRequestPayload{}
	err := json.Unmarshal([]byte(param), p)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	collection, cerr := store.GetCollection(p.Name)
	if cerr != nil {
		return "", fmt.Errorf("collection getting error: %w", err)
	}

	doc, derr := collection.Get(p.Id)
	if derr != nil {
		return "", fmt.Errorf("document getting error: %w", err)
	}

	r := GetDocCommandNameResponsePayload{
		Status: "success",
		Value:  p.Id,
		Doc:    *doc,
	}

	resp, merr := json.Marshal(r)
	if merr != nil {
		log.Println("internal error:", merr)
		return "", errors.New("internal error")
	}
	return string(resp), nil
}

func ExecDeleteDoc(param string) (string, error) {
	p := &DeleteDocCommandNameRequestPayload{}
	err := json.Unmarshal([]byte(param), p)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	collection, cerr := store.GetCollection(p.Name)
	if cerr != nil {
		return "", fmt.Errorf("collection getting error: %w", err)
	}

	derr := collection.Delete(p.Id)
	if derr != nil {
		return "", fmt.Errorf("document getting error: %w", err)
	}

	r := DeleteDocCommandNameResponsePayload{
		Status: "success",
		Value:  p.Id,
	}

	resp, merr := json.Marshal(r)
	if merr != nil {
		log.Println("internal error:", merr)
		return "", errors.New("internal error")
	}
	return string(resp), nil
}
