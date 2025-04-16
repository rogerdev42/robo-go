package documentstore

import "errors"

type Store struct {
	collections map[string]*Collection
}

var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
)

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if _, ok := s.collections[name]; ok {
		return nil, ErrCollectionAlreadyExists
	}
	s.collections[name] = &Collection{
		cfg:       *cfg,
		documents: make(map[string]Document),
	}
	return s.collections[name], nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	if collection, ok := s.collections[name]; ok {
		return collection, nil
	} else {
		return nil, ErrCollectionNotFound
	}
}

func (s *Store) DeleteCollection(name string) error {
	if _, ok := s.collections[name]; !ok {
		return ErrCollectionNotFound

	}
	delete(s.collections, name)
	return nil
}
