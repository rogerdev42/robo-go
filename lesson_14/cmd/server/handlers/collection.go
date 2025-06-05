package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionHandler struct {
	collection *mongo.Collection
}

// Request body structures for collections
type CreateCollectionReqBody struct {
	Name string `json:"name"`
}

type DeleteCollectionReqBody struct {
	Name string `json:"name"`
}

func NewCollectionHandler(collection *mongo.Collection) *CollectionHandler {
	return &CollectionHandler{
		collection: collection,
	}
}

func (h *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody CreateCollectionReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("json parsing error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate collection name
	if reqBody.Name == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Check if collection with this name already exists
	filter := map[string]interface{}{
		"name": reqBody.Name,
	}

	var existingCollection map[string]interface{}
	err := h.collection.FindOne(context.Background(), filter).Decode(&existingCollection)

	if err == nil {
		log.Printf("collection with name '%s' already exists", reqBody.Name)
		http.Error(w, `{"ok": false}`, http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		log.Printf("error searching for collection: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	document := map[string]interface{}{
		"name":       reqBody.Name,
		"created_at": time.Now(),
	}

	result, err := h.collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Printf("mongoDB insertion error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	log.Printf("document successfully inserted with ID: %v", result.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok": true,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *CollectionHandler) ListCollections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	cursor, err := h.collection.Find(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("error getting collections list: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var collections []map[string]interface{}
	if err = cursor.All(context.Background(), &collections); err != nil {
		log.Printf("error reading collections: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok":          true,
		"count":       len(collections),
		"collections": collections,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *CollectionHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody DeleteCollectionReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("json parsing error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate collection name
	if reqBody.Name == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	filter := map[string]interface{}{
		"name": reqBody.Name,
	}

	result, err := h.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("mongoDB deletion error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("collection with name '%s' not found", reqBody.Name)
		http.Error(w, `{"ok": false}`, http.StatusNotFound)
		return
	}

	log.Printf("Ccllection '%s' successfully deleted", reqBody.Name)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok": true,
	}

	json.NewEncoder(w).Encode(response)
}
