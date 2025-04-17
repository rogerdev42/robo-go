package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type Store struct {
	collections map[string]*Collection
}

var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
	ErrDumpStore               = errors.New("failed to dump store")
	ErrDumpStoreFile           = errors.New("failed to write dump to file")
	ErrReadFile                = errors.New("failed to read file")
)

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if _, ok := s.collections[name]; ok {
		l.Error("collection creation error: collection already exists", slog.Any("name", name))
		return nil, ErrCollectionAlreadyExists
	}
	s.collections[name] = &Collection{
		cfg:       *cfg,
		documents: make(map[string]Document),
	}
	l.Info("collection created", slog.Any("name", name))
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
		l.Error("collection deletion error: collection not found", slog.Any("name", name))
		return ErrCollectionNotFound
	}
	delete(s.collections, name)
	return nil
}

func (s *Store) Dump() ([]byte, error) {

	l.Info("hello dumping store", slog.Attr{
		Key:   "qwe",
		Value: slog.AnyValue(56),
	})
	dump := map[string]any{
		"collections": make(map[string]any),
	}

	for name, collection := range s.collections {
		collectionDump := map[string]any{
			"config":    collection.cfg,
			"documents": collection.documents,
		}
		dump["collections"].(map[string]any)[name] = collectionDump
	}

	jsonDump, err := json.Marshal(dump)
	if err == nil {
		return jsonDump, nil
	}

	return nil, fmt.Errorf(ErrDumpStore.Error()+": %w", err)
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	var data struct {
		Collections map[string]struct {
			Config    CollectionConfig    `json:"config"`
			Documents map[string]Document `json:"documents"`
		} `json:"collections"`
	}

	err := json.Unmarshal(dump, &data)
	if err != nil {
		return nil, fmt.Errorf(ErrDumpStore.Error()+": %w", err)
	}

	store := &Store{
		collections: make(map[string]*Collection),
	}

	for name, data := range data.Collections {
		collection := &Collection{
			cfg:       data.Config,
			documents: data.Documents,
		}
		store.collections[name] = collection
	}

	return store, nil
}

func (s *Store) DumpToFile(filename string) error {
	data, err := s.Dump()
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf(ErrDumpStoreFile.Error()+": %w", err)
	}

	return nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf(ErrReadFile.Error()+": %w", err)
	}

	store, err := NewStoreFromDump(data)
	if err != nil {
		return nil, err
	}

	return store, nil
}
