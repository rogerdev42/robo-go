package lru

import (
	"log/slog"
)

type LruCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type LRUCache struct {
	capacity int
	cache    map[string]*node
	list     *List
}

func NewLruCache(capacity int) LruCache {
	if capacity <= 0 {
		return nil
	}
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*node),
		list:     NewList(),
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	node, ok := c.cache[key]
	if !ok {
		return "", false
	}
	if err := c.list.MoveToFront(node); err != nil {
		slog.Error(err.Error())
		return "", false
	}

	return node.value, true
}

func (c *LRUCache) Put(key, value string) {
	if node, ok := c.cache[key]; ok {
		node.value = value
		if err := c.list.MoveToFront(node); err != nil {
			slog.Error(err.Error())
		}
		return
	}

	newNode := c.list.PushFront(value)
	c.cache[key] = newNode

	if c.list.Len() > c.capacity {
		lastNode := c.list.Back()
		if lastNode != nil {
			for k, n := range c.cache {
				if n == lastNode {
					delete(c.cache, k)
					break
				}
			}
			if err := c.list.Remove(lastNode); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
