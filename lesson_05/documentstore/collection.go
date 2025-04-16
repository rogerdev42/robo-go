package documentstore

import "errors"

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

var (
	ErrDocumentNotFound        = errors.New("document not found")
	ErrDocumentNoPrimaryKey    = errors.New("document must have a primary key of type string")
	ErrDocumentInvalidKeyType  = errors.New("primary key value must be a string")
	ErrDocumentEmptyPrimaryKey = errors.New("primary key value cannot be empty\n")
)

func (s *Collection) Put(doc Document) error {
	key, ok := doc.Fields[s.cfg.PrimaryKey]
	if !ok || key.Type != DocumentFieldTypeString {
		return ErrDocumentNoPrimaryKey
	}

	pk, ok := key.Value.(string)
	if !ok {
		return ErrDocumentInvalidKeyType
	}
	if pk == "" {
		return ErrDocumentEmptyPrimaryKey
	}

	s.documents[pk] = doc
	return nil
}

func (s *Collection) Get(key string) (*Document, error) {
	if doc, ok := s.documents[key]; ok {
		return &doc, nil
	}
	return nil, ErrDocumentNotFound
}

func (s *Collection) Delete(key string) error {
	if _, ok := s.documents[key]; !ok {
		return ErrDocumentNotFound
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
