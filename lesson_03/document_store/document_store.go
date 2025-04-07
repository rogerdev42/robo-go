package document_store

import "fmt"

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]Document{}

func Put(doc Document) bool {
	if err := ValidateDocument(doc); err != nil {
		println(err.Error())
		return false
	}

	key := doc.Fields["key"].Value.(string)
	documents[key] = doc
	return true
}

func ValidateDocument(doc Document) error {
	field, ok := doc.Fields["key"]
	if !ok {
		return fmt.Errorf("document must have a key field")
	}

	if field.Type != DocumentFieldTypeString {
		return fmt.Errorf("key field value must be a string")
	}

	return nil
}

func Get(key string) (*Document, bool) {
	doc, ok := documents[key]
	if !ok {
		return nil, false
	}

	return &doc, true
}

func Delete(key string) bool {
	if _, ok := documents[key]; !ok {
		return false
	}

	delete(documents, key)
	return true

}

func List() []Document {
	docs := make([]Document, 0, len(documents))
	for _, doc := range documents {
		docs = append(docs, doc)
	}
	return docs
}
