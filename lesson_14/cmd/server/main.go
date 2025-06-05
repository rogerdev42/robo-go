package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"lesson_14/cmd/server/handlers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	export MONGO_HOST="localhost"
	export MONGO_PORT="27017"
	export MONGO_USERNAME="root"
	export MONGO_PASSWORD="root"
	export DB_NAME="app"
	export DOCUMENTS_COLLECTION="documents"
	export COLLECTIONS_COLLECTION="collections"
	export PORT="8080"
*/

func main() {
	// Get configuration from environment variables
	mongoURI := buildMongoURI()
	dbName := getEnv("DB_NAME", "app")
	documentsCollectionName := getEnv("DOCUMENTS_COLLECTION", "documents")
	collectionsCollectionName := getEnv("COLLECTIONS_COLLECTION", "collection")
	port := getEnv("PORT", "8080")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("mongoDB connection error:", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("error disconnecting from MongoDB:", err)
		}
	}(client, ctx)

	documentsCollection := client.Database(dbName).Collection(documentsCollectionName)
	collectionsCollection := client.Database(dbName).Collection(collectionsCollectionName)
	documentHandler := handlers.NewDocumentHandler(documentsCollection)
	collectionHandler := handlers.NewCollectionHandler(collectionsCollection)

	http.HandleFunc("/put_document", documentHandler.PutDocument)
	http.HandleFunc("/get_document", documentHandler.GetDocument)
	http.HandleFunc("/list_documents", documentHandler.ListDocuments)
	http.HandleFunc("/delete_document", documentHandler.DeleteDocument)

	http.HandleFunc("/create_collection", collectionHandler.CreateCollection)
	http.HandleFunc("/list_collections", collectionHandler.ListCollections)
	http.HandleFunc("/delete_collection", collectionHandler.DeleteCollection)

	fmt.Printf("server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildMongoURI() string {
	// If full URI is provided, use it
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}

	// Otherwise build from components
	host := getEnv("MONGO_HOST", "localhost")
	port := getEnv("MONGO_PORT", "27017")
	username := getEnv("MONGO_USERNAME", "")
	password := getEnv("MONGO_PASSWORD", "")

	if username != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}
