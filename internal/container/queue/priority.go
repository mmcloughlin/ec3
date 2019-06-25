// Package queue implements queue types.
package queue

import "container/heap"

// Priority implements a prioritized queue. Low values are high priority (min heap).
type Priority struct {
	h *entryheap
}

// NewPriority creates an empty priority queue.
func NewPriority() *Priority {
	return &Priority{
		h: &entryheap{},
	}
}

// Empty reports whether the queue is empty.
func (q *Priority) Empty() bool {
	return q.h.Len() == 0
}

// Insert an item with the given priority.
func (q *Priority) Insert(i interface{}, p float64) {
	heap.Push(q.h, entry{item: i, priority: p})
}

// Peek at the highest priority item.
func (q *Priority) Peek() interface{} {
	return q.h.entries[0].item
}

// Pop off the highest priority item.
func (q *Priority) Pop() interface{} {
	x := heap.Pop(q.h)
	return x.(entry).item
}

// entry is an entry in the priority queue.
type entry struct {
	item     interface{}
	priority float64
}

// entryheap implements a heap of entries.
type entryheap struct {
	entries []entry
}

func (h entryheap) Len() int {
	return len(h.entries)
}

func (h entryheap) Less(i, j int) bool {
	return h.entries[i].priority < h.entries[j].priority
}

func (h entryheap) Swap(i, j int) {
	h.entries[i], h.entries[j] = h.entries[j], h.entries[i]
}

func (h *entryheap) Push(x interface{}) {
	h.entries = append(h.entries, x.(entry))
}

func (h *entryheap) Pop() interface{} {
	n := len(h.entries)
	x := h.entries[n-1]
	h.entries = h.entries[:n-1]
	return x
}
