package documentstore

type Collection struct {
	cfg       CollectionConfig
	Documents map[string]Document
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
	s.Documents[key.Value.(string)] = doc
}

func (s *Collection) Get(key string) (*Document, bool) {
	if doc, ok := s.Documents[key]; ok {
		return &doc, true
	} else {
		return nil, false
	}
}

func (s *Collection) Delete(key string) bool {
	if _, ok := s.Documents[key]; !ok {
		return false
	}
	delete(s.Documents, key)
	return true
}

func (s *Collection) List() []Document {
	documents := make([]Document, 0, len(s.Documents))
	for _, doc := range s.Documents {
		documents = append(documents, doc)
	}
	return documents
}
