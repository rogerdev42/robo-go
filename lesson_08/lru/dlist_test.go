package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	list := NewList()
	assert.NotNil(t, list, "NewList should return non-nil list")
	assert.Equal(t, 0, list.Len(), "New list should be empty")
	assert.NotNil(t, list.head, "Head should be initialized")
	assert.NotNil(t, list.tail, "Tail should be initialized")
	assert.Equal(t, list.tail, list.head.next, "Head should point to tail")
	assert.Equal(t, list.head, list.tail.prev, "Tail should point to head")
}

func TestPushFront(t *testing.T) {
	list := NewList()

	t.Run("Add first node", func(t *testing.T) {
		node1 := list.PushFront("hello")
		assert.Equal(t, 1, list.Len(), "Len should be 1 after PushFront")
		assert.Equal(t, "hello", node1.value, "Node value should be hello")
		assert.Equal(t, list.head, node1.prev, "Node prev should be head")
		assert.Equal(t, list.tail, node1.next, "Node next should be tail")
		assert.Equal(t, node1, list.head.next, "Head should point to node1")
		assert.Equal(t, node1, list.tail.prev, "Tail should point to node1")
	})

	t.Run("Add second node", func(t *testing.T) {
		node2 := list.PushFront("world")
		assert.Equal(t, 2, list.Len(), "Len should be 2 after second PushFront")
		assert.Equal(t, "world", node2.value, "Node value should be world")
		assert.Equal(t, list.head, node2.prev, "Node2 prev should be head")
		assert.Equal(t, list.head.next.next, node2.next, "Node2 next should be node1")
		assert.Equal(t, node2, list.head.next.next.prev, "Node1 prev should be node2")
		assert.Equal(t, list.tail, list.head.next.next.next, "Node1 next should be tail")
		assert.Equal(t, node2, list.head.next, "Head should point to node2")
	})
}

func TestRemove(t *testing.T) {
	list := NewList()
	node := list.PushFront("hello")

	t.Run("Remove valid node", func(t *testing.T) {
		err := list.Remove(node)
		assert.NoError(t, err, "Remove valid node should not return error")
		assert.Equal(t, 0, list.Len(), "Len should be 0 after Remove")
		assert.Equal(t, list.tail, list.head.next, "Head should point to tail")
		assert.Equal(t, list.head, list.tail.prev, "Tail should point to head")
		assert.Nil(t, node.prev, "Removed node prev should be nil")
		assert.Nil(t, node.next, "Removed node next should be nil")
	})

	t.Run("Remove invalid node", func(t *testing.T) {
		err := list.Remove(node)
		assert.Error(t, err, "Remove invalid node should return error")
		assert.Equal(t, "invalid node", err.Error(), "Expected invalid node error")
		err = list.Remove(nil)
		assert.Error(t, err, "Remove nil node should return error")
		assert.Equal(t, "invalid node", err.Error(), "Expected invalid node error")
	})
}

func TestMoveToFront(t *testing.T) {
	list := NewList()
	node1 := list.PushFront("hello")
	node2 := list.PushFront("world")

	//t.Run("Move node already at front", func(t *testing.T) {
	//	err := list.MoveToFront(node2)
	//	assert.Error(t, err, "MoveToFront node at front should return error")
	//	assert.Equal(t, "node already at front", err.Error(), "Expected node already at front error")
	//	assert.Equal(t, node2, list.head.next, "Node2 should still be at front")
	//})

	t.Run("Move valid node to front", func(t *testing.T) {
		err := list.MoveToFront(node1)
		assert.NoError(t, err, "MoveToFront valid node should not return error")
		assert.Equal(t, node1, list.head.next, "Node1 should be at front")
		assert.Equal(t, node2, node1.next, "Node1 next should be node2")
		assert.Equal(t, list.head, node1.prev, "Node1 prev should be head")
		assert.Equal(t, node1, node2.prev, "Node2 prev should be node1")
		assert.Equal(t, list.tail, node2.next, "Node2 next should be tail")
	})

	t.Run("Move invalid node", func(t *testing.T) {
		err := list.Remove(node1)
		assert.NoError(t, err, "Remove node1 should succeed")
		err = list.MoveToFront(node1)
		assert.Error(t, err, "MoveToFront invalid node should return error")
		assert.Equal(t, "invalid node", err.Error(), "Expected invalid node error")
		err = list.MoveToFront(nil)
		assert.Error(t, err, "MoveToFront nil node should return error")
		assert.Equal(t, "invalid node", err.Error(), "Expected invalid node error")
	})
}

func TestBack(t *testing.T) {
	list := NewList()

	t.Run("Empty list", func(t *testing.T) {
		assert.Nil(t, list.Back(), "Back on empty list should return nil")
	})

	t.Run("Single node", func(t *testing.T) {
		node1 := list.PushFront("hello")
		back := list.Back()
		assert.Equal(t, node1, back, "Back should return node1")
		assert.Equal(t, "hello", back.value, "Back node value should be hello")
	})

	t.Run("Multiple nodes", func(t *testing.T) {
		list.PushFront("hello") // Первый узел
		list.PushFront("world") // Второй узел
		back := list.Back()
		assert.Equal(t, "hello", back.value, "Back node value should be hello")
	})
}

func TestLen(t *testing.T) {
	list := NewList()

	t.Run("Empty list", func(t *testing.T) {
		assert.Equal(t, 0, list.Len(), "Len of empty list should be 0")
	})

	t.Run("Add nodes", func(t *testing.T) {
		list.PushFront("hello")
		assert.Equal(t, 1, list.Len(), "Len should be 1 after PushFront")
		list.PushFront("world")
		assert.Equal(t, 2, list.Len(), "Len should be 2 after second PushFront")
	})

	t.Run("Remove node", func(t *testing.T) {
		err := list.Remove(list.head.next)
		assert.NoError(t, err, "Remove should succeed")
		assert.Equal(t, 1, list.Len(), "Len should be 1 after Remove")
	})
}
