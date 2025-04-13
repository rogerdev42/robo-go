package documentstore

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) {
	key, ok := doc.Fields[s.cfg.PrimaryKey]
	if !ok || key.Type != DocumentFieldTypeString {
		println("Document must have a primary key field of type string")
		return
	}

	pk, ok := key.Value.(string)
	if !ok {
		println("Primary key value must be a string")
		return
	}
	if pk == "" {
		println("Primary key value must not be empty")
		return
	}

	s.documents[pk] = doc
}

func (s *Collection) Get(key string) (*Document, bool) {
	if doc, ok := s.documents[key]; ok {
		return &doc, true
	}
	return nil, false
}

func (s *Collection) Delete(key string) bool {
	if _, ok := s.documents[key]; !ok {
		return false
	}
	delete(s.documents, key)
	return true
}

func (s *Collection) List() []Document {
	documents := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		documents = append(documents, doc)
	}
	return documents
}
