package documentstore

import (
	"errors"
	"lesson_09/pkg/bst"
	"log/slog"
)

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
	indexes   map[string]*bst.BinarySearchTree
}

type CollectionConfig struct {
	PrimaryKey string
}

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

var (
	ErrDocumentNotFound        = errors.New("document not found")
	ErrDocumentNoPrimaryKey    = errors.New("document must have a primary key of type string")
	ErrDocumentInvalidKeyType  = errors.New("primary key value must be a string")
	ErrDocumentEmptyPrimaryKey = errors.New("primary key value cannot be empty")
)

func (s *Collection) Put(doc Document) error {
	key, ok := doc.Fields[s.cfg.PrimaryKey]
	if !ok || key.Type != DocumentFieldTypeString {
		l.Error("document creation error: document with no primary key", slog.Any("document", doc))
		return ErrDocumentNoPrimaryKey
	}

	pk, ok := key.Value.(string)
	if !ok {
		l.Error("document creation error: document with invalid primary key", slog.Any("PrimaryKey", key))
		return ErrDocumentInvalidKeyType
	}
	if pk == "" {
		l.Error("document creation error: document with empty key", slog.Any("document", doc))
		return ErrDocumentEmptyPrimaryKey
	}

	// If a document already exists, delete it from the index before updating it.
	if oldDoc, ok := s.documents[pk]; ok {
		for fieldName, index := range s.indexes {
			field := oldDoc.Fields[fieldName]
			if field.Type == DocumentFieldTypeString {
				if strValue, ok := field.Value.(string); ok {
					err := index.Delete(strValue)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Update or insert the document
	s.documents[pk] = doc

	// Add
	for fieldName, index := range s.indexes {
		field := doc.Fields[fieldName]
		if field.Type == DocumentFieldTypeString {
			if strValue, ok := field.Value.(string); ok {
				err := index.Insert(strValue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Collection) Get(key string) (*Document, error) {
	if doc, ok := s.documents[key]; ok {
		return &doc, nil
	}
	return nil, ErrDocumentNotFound
}

func (s *Collection) Delete(key string) error {
	doc, ok := s.documents[key]
	if !ok {
		l.Error("document deletion error: document not found", slog.Any("PrimaryKey", key))
		return ErrDocumentNotFound
	}

	for fieldName, index := range s.indexes {
		field := doc.Fields[fieldName]
		if field.Type == DocumentFieldTypeString {
			if strValue, ok := field.Value.(string); ok {
				err := index.Delete(strValue)
				if err != nil {
					return err
				}
			}
		}
	}

	delete(s.documents, key)

	return nil
}

func (s *Collection) List() []Document {
	documents := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		documents = append(documents, doc)
	}
	return documents
}

func (s *Collection) CreateIndex(fieldName string) error {
	if _, ok := s.indexes[fieldName]; ok {
		return errors.New("index already exists")
	}
	tree := bst.NewBST()
	for _, doc := range s.documents {
		field := doc.Fields[fieldName]
		if field.Type == DocumentFieldTypeString {
			if strValue, ok := field.Value.(string); ok {
				err := tree.Insert(strValue)
				if err != nil {
					return err
				}
			}
		}
	}
	s.indexes[fieldName] = tree

	return nil
}

func (s *Collection) DeleteIndex(fieldName string) error {
	if _, ok := s.indexes[fieldName]; !ok {
		return errors.New("index not found")
	}
	delete(s.indexes, fieldName)
	return nil
}

func (s *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	index, ok := s.indexes[fieldName]
	if !ok {
		return nil, errors.New("index not found")
	}

	keys := index.RangeTraversal(*params.MinValue, *params.MaxValue, params.Desc)

	result := make([]Document, 0, len(keys))
	for _, key := range keys {
		for _, doc := range s.documents {
			field, exists := doc.Fields[fieldName]
			if !exists {
				continue
			}

			if field.Type != DocumentFieldTypeString {
				continue
			}

			if fieldValue, ok := field.Value.(string); ok && fieldValue == key {
				result = append(result, doc)
			}
		}

	}
	return result, nil
}
