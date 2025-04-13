package documentstore

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	if _, ok := s.collections[name]; ok {
		println("collection already exists")
		return false, nil
	}
	s.collections[name] = &Collection{
		cfg:       *cfg,
		documents: make(map[string]Document),
	}
	return true, s.collections[name]
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	if collection, ok := s.collections[name]; ok {
		return collection, true
	} else {
		println("collection does not exist")
		return nil, false
	}
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.collections[name]; !ok {
		println("collection does not exist")
		return false

	}
	delete(s.collections, name)
	return true
}
