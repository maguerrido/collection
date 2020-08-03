// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package queue implements a singly-linked list with queue behaviors.
package queue

import (
	"fmt"
	coll "github.com/maguerrido/collection"
)

// node of a Queue.
type node struct {
	// value stored in the node.
	value interface{}

	// next points to the next node.
	// If the node is the back node, then points to nil.
	next *node
}

// clear sets the properties of the node to its zero values.
// Time complexity: O(1).
func (n *node) clear() {
	n.value, n.next = nil, nil
}

// Queue represents a singly-linked list.
// The zero value for Queue is an empty Queue ready to use.
type Queue struct {
	// front points to the front (first) node in the queue.
	// back points to the back (last) node in the queue.
	front, back *node

	// len is the current length (number of nodes).
	len int
}

// New returns a new Queue ready to use.
// Time complexity: O(1).
func New() *Queue {
	return new(Queue)
}

// NewBySlice returns a new Queue with the values stored in the slice keeping its order.
// Time complexity: O(n), where n is the current length of the slice.
func NewBySlice(values []interface{}) *Queue {
	q := New()
	for _, v := range values {
		q.Push(v)
	}
	return q
}

// Clone returns a new cloned Queue.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) Clone() *Queue {
	clone := New()
	for n := q.front; n != nil; n = n.next {
		clone.Push(n.value)
	}
	return clone
}

// Do gets the front value and performs all the procedures, then repeats it with the rest of the values.
// The queue will be empty.
// Time complexity: O(n*p), where n is the current length of the queue and p is the number of procedures.
func (q *Queue) Do(procedures ...func(v interface{})) {
	for !q.IsEmpty() {
		v := q.Get()
		for _, procedure := range procedures {
			procedure(v)
		}
	}
}

// Equals compares this queue with the 'other' queue and returns true if they are equal.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) Equals(other *Queue) bool {
	if q.len != other.len {
		return false
	}
	for i, j := q.front, other.front; i != nil; i, j = i.next, j.next {
		if i.value != j.value {
			return false
		}
	}
	return true
}

// EqualsByComparator compares this queue with the 'other' queue and returns true if they are equal.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) EqualsByComparator(other *Queue, equals func(v1, v2 interface{}) bool) bool {
	if q.len != other.len {
		return false
	}
	for i, j := q.front, other.front; i != nil; i, j = i.next, j.next {
		if !equals(i.value, j.value) {
			return false
		}
	}
	return true
}

// Get returns the front value and removes it from the queue.
// If the queue is empty, then returns nil.
// Time complexity: O(1).
func (q *Queue) Get() interface{} {
	if q.IsEmpty() {
		return nil
	}
	n, v := q.front, q.front.value
	if q.back == n {
		q.back = nil
	}
	q.front = n.next
	n.clear()
	q.len--
	return v
}

// GetIf returns all first values that meet the condition defined by the 'condition' parameter. These values will be
//removed from the queue.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) GetIf(condition func(v interface{}) bool) []interface{} {
	values := make([]interface{}, 0)
	for n := q.front; n != nil; {
		if condition(n.value) {
			values = append(values, q.Get())
			n = q.front
		} else {
			return values
		}
	}
	return values
}

// IsEmpty returns true if the queue has no values.
// Time complexity: O(1).
func (q *Queue) IsEmpty() bool {
	return q.len == 0
}

func (q *Queue) Iterator() coll.Iterator {
	return &iterator{
		q:           q,
		prev:        nil,
		this:        nil,
		index:       -1,
		lastCommand: -1,
		lastHasNext: false,
	}
}

// Len returns the current length of the queue.
// Time complexity: O(1).
func (q *Queue) Len() int {
	return q.len
}

// Peek returns the front value.
// If the queue is empty, then returns nil.
// Time complexity: O(1).
func (q *Queue) Peek() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.front.value
}

// Push inserts the value 'v' at the back of the queue.
// Time complexity: O(1).
func (q *Queue) Push(v interface{}) {
	n := &node{value: v, next: nil}
	if q.IsEmpty() {
		q.front = n
	} else {
		q.back.next = n
	}
	q.back = n
	q.len++
}

// RemoveAll sets the properties of the queue to its zero values.
// Time complexity: O(1).
func (q *Queue) RemoveAll() {
	q.front, q.back, q.len = nil, nil, 0
}

func (q *Queue) removeNode(prev, n *node) {
	if q.front == n {
		if q.back == n {
			q.front, q.back = nil, nil
		} else {
			q.front = n.next
		}
	} else if q.back == n {
		q.back = prev
		prev.next = nil
	} else {
		prev.next = n.next
	}
	n.clear()
	q.len--
}

// Search returns the index (zero based) of the first match of the value 'v'.
// If the value 'v' does not belong to the queue, then returns -1.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) Search(v interface{}) int {
	for n, i := q.front, 0; n != nil; n, i = n.next, i+1 {
		if n.value == v {
			return i
		}
	}
	return -1
}

// SearchByComparator returns the index (zero based) of the first match of the value 'v'.
// If the value 'v' does not belong to the queue, then returns -1.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) SearchByComparator(v interface{}, equals func(v1, v2 interface{}) bool) int {
	for n, i := q.front, 0; n != nil; n, i = n.next, i+1 {
		if equals(n.value, v) {
			return i
		}
	}
	return -1
}

// Slice returns a new slice with the values stored in the queue keeping its order.
// The queue retains its original state.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) Slice() []interface{} {
	values := make([]interface{}, 0, q.len)
	for n := q.front; n != nil; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// String returns a representation of the queue as a string.
// Queue implements the fmt.Stringer interface.
// Time complexity: O(n), where n is the current length of the queue.
func (q *Queue) String() string {
	if q.IsEmpty() {
		return "[]"
	}
	n := q.front
	str := "["
	for ; n.next != nil; n = n.next {
		str += fmt.Sprintf("%v ", n.value)
	}
	return str + fmt.Sprintf("%v]", n.value)
}

type iterator struct {
	q           *Queue
	prev, this  *node
	index       int
	lastCommand int
	lastHasNext bool
}

const (
	iteratorCommandHasNext = 0
	iteratorCommandNext    = 1
	iteratorCommandRemove  = 2
)

func (i *iterator) ForEach(action func(v *interface{})) {
	if action != nil {
		for e := i.q.front; e != nil; e = e.next {
			action(&e.value)
		}
	}
}

func (i *iterator) HasNext() bool {
	i.lastCommand = iteratorCommandHasNext
	i.lastHasNext = i.index < i.q.len-1
	return i.lastHasNext
}

func (i *iterator) Next() (interface{}, error) {
	if i.lastCommand != iteratorCommandHasNext {
		return nil, fmt.Errorf(coll.ErrorIteratorNext)
	} else if !i.lastHasNext {
		return nil, fmt.Errorf(coll.ErrorIteratorHasNext)
	}

	if i.this == nil {
		i.this = i.q.front
	} else {
		i.prev = i.this
		i.this = i.this.next
	}
	i.index++
	i.lastCommand = iteratorCommandNext

	return i.this.value, nil
}

func (i *iterator) Remove() error {
	if !i.lastHasNext {
		return fmt.Errorf(coll.ErrorIteratorHasNext)
	} else if i.lastCommand != iteratorCommandNext {
		return fmt.Errorf(coll.ErrorIteratorRemove)
	}

	i.q.removeNode(i.prev, i.this)
	i.this = i.prev
	i.index--
	i.lastCommand = iteratorCommandRemove

	return nil
}
