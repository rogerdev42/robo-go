package lru

import (
	"lesson_08/internal/documentstore"
	"math"
	"time"
)

type LruCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type cache struct {
	capacity   int
	collection *documentstore.Collection
	readtime   map[string]int64
}

func NewLruCache(capacity int) LruCache {
	store := documentstore.NewStore()
	collection, _ := store.CreateCollection("cache", &documentstore.CollectionConfig{PrimaryKey: "key"})

	return cache{
		capacity:   capacity,
		collection: collection,
		readtime:   make(map[string]int64),
	}
}

func (c cache) Put(key, value string) {
	if c.capacity == len(c.collection.List()) {
		key := c.getKeyWithMinReadTime()
		c.collection.Delete(key)
		delete(c.readtime, key)
	}

	doc := documentstore.Document{
		Fields: make(map[string]documentstore.DocumentField),
	}
	doc.Fields["key"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeString,
		Value: key,
	}
	doc.Fields["value"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeString,
		Value: value,
	}

	c.readtime[key] = (time.Now()).UnixMicro()

	c.collection.Put(doc)
}

func (c cache) Get(key string) (string, bool) {
	doc, err := c.collection.Get(key)
	if err != nil {
		return "", false
	}
	field, ok := doc.Fields["value"]
	if !ok {
		return "", false
	}
	str, ok := field.Value.(string)
	if !ok {
		return "", false
	}

	c.readtime[key] = (time.Now()).UnixMicro()

	return str, true
}

func (c cache) getKeyWithMinReadTime() string {
	if len(c.readtime) == 0 {
		return ""
	}

	var minKey string
	var minTime int64 = math.MaxInt64

	for key, timestamp := range c.readtime {
		if timestamp < minTime {
			minTime = timestamp
			minKey = key
		}
	}

	return minKey
}
