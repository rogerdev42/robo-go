package lru

import "errors"

// node represents a node in the doubly linked list.
// This is an internal structure used only through List methods.
type node struct {
	value string
	prev  *node
	next  *node
}

// List represents a doubly linked list with sentinel nodes.
type List struct {
	head *node
	tail *node
	len  int
}

// NewList creates a new empty doubly linked list with sentinel nodes.
func NewList() *List {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head
	return &List{
		head: head,
		tail: tail,
	}
}

// PushFront adds a new node with the given value to the front of the list and returns a pointer to it.
func (l *List) PushFront(value string) *node {
	newNode := &node{value: value}
	newNode.next = l.head.next
	newNode.prev = l.head
	newNode.next.prev = newNode
	l.head.next = newNode
	l.len++
	return newNode
}

// Remove removes the specified node from the list.
// It returns an error if the node is invalid (nil or already removed).
func (l *List) Remove(node *node) error {
	if node == nil || node.prev == nil || node.next == nil {
		return errors.New("invalid node")
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev = nil
	node.next = nil
	l.len--
	return nil
}

// MoveToFront moves the specified node to the front of the list.
// It returns an error if the node is invalid or already at the front.
func (l *List) MoveToFront(node *node) error {
	if node == nil || node.prev == nil || node.next == nil {
		return errors.New("invalid node")
	}
	//if l.head.next == node {
	//	return errors.New("node already at front")
	//}
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = l.head.next
	node.prev = l.head
	l.head.next.prev = node
	l.head.next = node
	return nil
}

// Back returns the last node in the list or nil if the list is empty.
func (l *List) Back() *node {
	if l.len == 0 {
		return nil
	}
	return l.tail.prev
}

// Len returns the number of nodes in the list.
func (l *List) Len() int {
	return l.len
}
