package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentHandler struct {
	collection *mongo.Collection
}

// Request body structures for documents
type PutDocumentReqBody struct {
	CollectionName string                 `json:"collection_name"`
	Document       map[string]interface{} `json:"document"`
}

type GetDocumentReqBody struct {
	CollectionName string `json:"collection_name"`
	ID             string `json:"id"`
}

type ListDocumentsReqBody struct {
	CollectionName string `json:"collection_name"`
}

type DeleteDocumentReqBody struct {
	CollectionName string `json:"collection_name"`
	ID             string `json:"id"`
}

func NewDocumentHandler(collection *mongo.Collection) *DocumentHandler {
	return &DocumentHandler{
		collection: collection,
	}
}

func (h *DocumentHandler) PutDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody PutDocumentReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("json parsing error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reqBody.CollectionName == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	if reqBody.Document == nil || len(reqBody.Document) == 0 {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Check for required "id" field
	documentId, exists := reqBody.Document["id"]
	if !exists || documentId == nil {
		log.Printf("missing required 'id' field in document")
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	targetCollection := h.collection.Database().Collection(reqBody.CollectionName)

	filter := map[string]interface{}{
		"id": documentId,
	}

	// Use ReplaceOne with upsert=true for create or replace
	opts := options.ReplaceOptions{}
	opts.SetUpsert(true)

	result, err := targetCollection.ReplaceOne(context.Background(), filter, reqBody.Document, &opts)
	if err != nil {
		log.Printf("document create/update error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	if result.UpsertedCount > 0 {
		log.Printf("new document created with ID: %v", documentId)
	} else if result.ModifiedCount > 0 {
		log.Printf("document with ID %v successfully updated", documentId)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok": true,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody GetDocumentReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("json parsing error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reqBody.CollectionName == "" || reqBody.ID == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	targetCollection := h.collection.Database().Collection(reqBody.CollectionName)

	filter := map[string]interface{}{
		"id": reqBody.ID,
	}

	var document map[string]interface{}
	err := targetCollection.FindOne(context.Background(), filter).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("document with ID %s not found", reqBody.ID)
			http.Error(w, `{"ok": false}`, http.StatusNotFound)
			return
		}
		log.Printf("document search error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok":       true,
		"document": document,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *DocumentHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody ListDocumentsReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("json parsing error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reqBody.CollectionName == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	targetCollection := h.collection.Database().Collection(reqBody.CollectionName)

	cursor, err := targetCollection.Find(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("error getting documents list: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var documents []map[string]interface{}
	if err = cursor.All(context.Background(), &documents); err != nil {
		log.Printf("error reading documents: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok":        true,
		"count":     len(documents),
		"documents": documents,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody DeleteDocumentReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("json parsing error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reqBody.CollectionName == "" || reqBody.ID == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	targetCollection := h.collection.Database().Collection(reqBody.CollectionName)

	filter := map[string]interface{}{
		"id": reqBody.ID,
	}

	result, err := targetCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("document deletion error: %v", err)
		http.Error(w, `{"ok": false}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("document with ID %s not found", reqBody.ID)
		http.Error(w, `{"ok": false}`, http.StatusNotFound)
		return
	}

	log.Printf("document with ID %s successfully deleted", reqBody.ID)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"ok": true,
	}

	json.NewEncoder(w).Encode(response)
}
