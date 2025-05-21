package documentstore

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/google/btree"
)

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
	indexes   map[string]*btree.BTreeG[*indexItem]
	mu        sync.RWMutex
}

type CollectionConfig struct {
	PrimaryKey string
}

var (
	ErrDocumentNotFound        = errors.New("document not found")
	ErrDocumentNoPrimaryKey    = errors.New("document must have a primary key of type string")
	ErrDocumentInvalidKeyType  = errors.New("primary key value must be a string")
	ErrDocumentEmptyPrimaryKey = errors.New("primary key value cannot be empty")
	ErrIndexExists             = errors.New("index already exists")
	ErrIndexNotFound           = errors.New("index not found")
)

type indexItem struct {
	key string
	pk  string // primary key
}

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

	s.mu.Lock()
	defer s.mu.Unlock()

	oldDoc, existed := s.documents[pk]
	if existed && s.indexes != nil {
		for field, tree := range s.indexes {
			if oldField, ok := oldDoc.Fields[field]; ok && oldField.Type == DocumentFieldTypeString {
				if strVal, ok := oldField.Value.(string); ok {
					tree.Delete(&indexItem{key: strVal, pk: pk})
				}
			}
		}
	}

	s.documents[pk] = doc

	if s.indexes != nil {
		for field, tree := range s.indexes {
			if newField, ok := doc.Fields[field]; ok && newField.Type == DocumentFieldTypeString {
				if strVal, ok := newField.Value.(string); ok {
					tree.ReplaceOrInsert(&indexItem{key: strVal, pk: pk})
				}
			}
		}
	}

	l.Info("document created", slog.Any("PrimaryKey", key))
	return nil
}

func (s *Collection) Get(key string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if doc, ok := s.documents[key]; ok {
		return &doc, nil
	}
	return nil, ErrDocumentNotFound
}

func (s *Collection) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.documents[key]
	if !ok {
		l.Error("document deletion error: document not found", slog.Any("PrimaryKey", key))
		return ErrDocumentNotFound
	}
	if s.indexes != nil {
		for field, tree := range s.indexes {
			if fieldVal, ok := doc.Fields[field]; ok && fieldVal.Type == DocumentFieldTypeString {
				if strVal, ok := fieldVal.Value.(string); ok {
					tree.Delete(&indexItem{key: strVal, pk: key})
				}
			}
		}
	}
	delete(s.documents, key)
	return nil
}

func (s *Collection) List() []Document {
	s.mu.RLock()
	defer s.mu.RUnlock()

	documents := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		documents = append(documents, doc)
	}
	return documents
}

func (a *indexItem) Less(b *indexItem) bool {
	return a.key < b.key || (a.key == b.key && a.pk < b.pk)
}

func (s *Collection) CreateIndex(fieldName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.indexes == nil {
		s.indexes = make(map[string]*btree.BTreeG[*indexItem])
	}
	if _, exists := s.indexes[fieldName]; exists {
		return ErrIndexExists
	}
	tree := btree.NewG(8, func(a, b *indexItem) bool { return a.Less(b) })
	for pk, doc := range s.documents {
		field, ok := doc.Fields[fieldName]
		if !ok || field.Type != DocumentFieldTypeString {
			continue
		}
		strVal, ok := field.Value.(string)
		if !ok {
			continue
		}
		tree.ReplaceOrInsert(&indexItem{key: strVal, pk: pk})
	}
	s.indexes[fieldName] = tree
	return nil
}

func (s *Collection) DeleteIndex(fieldName string) error {
	if s.indexes == nil {
		return ErrIndexNotFound
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}
	delete(s.indexes, fieldName)
	return nil
}

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

func (s *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.indexes == nil {
		return nil, ErrIndexNotFound
	}
	tree, ok := s.indexes[fieldName]
	if !ok {
		return nil, ErrIndexNotFound
	}
	var result []Document
	visit := func(item *indexItem) bool {
		doc := s.documents[item.pk]
		result = append(result, doc)
		return true
	}

	var minValue, maxValue string
	if params.MinValue != nil {
		minValue = *params.MinValue
	}
	if params.MaxValue != nil {
		maxValue = *params.MaxValue
	}
	if params.Desc {
		tree.Descend(func(item *indexItem) bool {
			if params.MinValue != nil && item.key < minValue {
				return false
			}
			if params.MaxValue != nil && item.key > maxValue {
				return true
			}
			return visit(item)
		})
	} else {
		tree.Ascend(func(item *indexItem) bool {
			if params.MinValue != nil && item.key < minValue {
				return true
			}
			if params.MaxValue != nil && item.key > maxValue {
				return false
			}
			return visit(item)
		})
	}
	return result, nil
}
