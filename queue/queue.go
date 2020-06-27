// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package queue implements a queue using a singly-linked list.
package queue

import "fmt"

// node of a Queue.
type node struct {
	// value stored in this node.
	value interface{}

	// next points to the next node.
	// If this node is the back node then next points to nil.
	next *node
}

// clear sets the properties of this node to its zero values.
func (n *node) clear() {
	n.value, n.next = nil, nil
}

// Queue represents a singly-linked list.
// The zero value for Queue is an empty Queue ready to use.
type Queue struct {
	// front points to the front (first) element on the list.
	// back points to the back (last) element on the list.
	front, back *node

	// len is the current length.
	len int
}

// New returns a new Queue ready to use.
func New() *Queue {
	return new(Queue)
}

// NewBySlice returns a new Queue with the values stored in the slice keeping its order.
// Time complexity: O(n), where n is the current length of 'values'.
func NewBySlice(values []interface{}) *Queue {
	q := New()
	for _, v := range values {
		q.Push(v)
	}
	return q
}

// Clone returns a new cloned Queue of 'q'.
// 'q' retains its original state.
// Time complexity: O(n), where n is the current length of 'q'.
func (q *Queue) Clone() *Queue {
	clone := New()
	for n := q.front; n != nil; n = n.next {
		clone.Push(n.value)
	}
	return clone
}

// Do gets the front value and performs all the procedures, then repeats this with the rest of the values.
// 'q' will be empty.
// Time complexity: O(n), where n is the current length of 'q'.
func (q *Queue) Do(procedures ...func(v interface{})) {
	for !q.IsEmpty() {
		v := q.Get()
		for _, procedure := range procedures {
			procedure(v)
		}
	}
}

// Equals compares queue 'q' with queue 'other' and returns true if they are equals.
// Time complexity: O(n), where n is the current length of 'q' as well as 'other'.
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

// EqualsByComparator compares queue 'q' with queue 'other' and returns true if they are equals.
// The comparison between values is defined by parameter 'equals'.
// 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of 'q' as well as 'other'.
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
// If the queue is empty returns nil.
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

// GetIf returns all first values that meet the condition defined by the "condition" parameter. These values
//will be removed from the queue.
// If the queue is empty returns an empty slice.
// Time complexity: O(n), where n is the current length of 'q'.
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

// IsEmpty returns true if the queue has no elements.
// Time complexity: O(1).
func (q *Queue) IsEmpty() bool {
	return q.len == 0
}

// Len returns the current length.
// Time complexity: O(1).
func (q *Queue) Len() int {
	return q.len
}

// Peek returns the front value.
// If the queue is empty returns nil.
// Time complexity: O(1).
func (q *Queue) Peek() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.front.value
}

// Push inserts the value 'v' in the back.
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

// RemoveAll sets the properties of this queue to its zero values.
// Time complexity: O(1).
func (q *Queue) RemoveAll() {
	q.front, q.back, q.len = nil, nil, 0
}

// Search returns returns the index (based on zero) of the first match of the value 'v'.
// If the value 'v' does not belong to the queue, return -1.
// Time complexity: O(n), where n is the current length of 'q'.
func (q *Queue) Search(v interface{}) int {
	for n, i := q.front, 0; n != nil; n, i = n.next, i+1 {
		if n.value == v {
			return i
		}
	}
	return -1
}

// SearchByComparator returns the index (based on zero) of the first match of the value 'v'.
// If the value 'v' does not belong to the queue, return -1.
// The comparison between values is defined by parameter 'equals'.
// 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of 'q'.
func (q *Queue) SearchByComparator(v interface{}, equals func(v1, v2 interface{}) bool) int {
	for n, i := q.front, 0; n != nil; n, i = n.next, i+1 {
		if equals(n.value, v) {
			return i
		}
	}
	return -1
}

// Slice returns a new slice with the values stored in the queue keeping its order.
// 'q' retains its original state.
// Time complexity: O(n), where n is the current length of 'q'.
func (q *Queue) Slice() []interface{} {
	values := make([]interface{}, 0, q.len)
	for n := q.front; n != nil; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// String returns a representation of the queue as a string.
// Queue implements the fmt.Stringer interface.
// Time complexity: O(n), where n is the current length of 'q'.
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
