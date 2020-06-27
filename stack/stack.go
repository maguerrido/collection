// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package stack implements a stack using a singly-linked list.
package stack

import "fmt"

// node of a Stack.
type node struct {
	// value stored in this node.
	value interface{}

	// next points to the next node.
	// If this node is the bottom element then next points to nil.
	next *node
}

// clear sets the properties of this node to its zero values.
func (n *node) clear() {
	n.value, n.next = nil, nil
}

// Stack represents a singly-linked list.
// The zero value for Stack is an empty Stack ready to use.
type Stack struct {
	// top points to the top node in the Stack.
	top *node

	// len is the current length.
	len int
}

// New returns a new Stack ready to use.
func New() *Stack {
	return new(Stack)
}

// NewBySlice returns a new Stack with the values stored in the slice keeping its order.
// The last value of the slice will be the top value of the stack.
// Time complexity: O(n), where n is the current length of 'values'.
func NewBySlice(values []interface{}) *Stack {
	s := New()
	for _, v := range values {
		s.Push(v)
	}
	return s
}

// Clone returns a new cloned Stack of 's'.
// 's' retains its original state.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) Clone() *Stack {
	return &Stack{cloneRecursive(s.top), s.len}
}

// quickSortRecursive is an auxiliary recursive function of the Stack Clone method.
func cloneRecursive(n *node) *node {
	if n == nil {
		return nil
	}
	return &node{n.value, cloneRecursive(n.next)}
}

// Do gets the top value and performs all the procedures, then repeats this with the rest of the values.
// 's' will be empty.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) Do(procedures ...func(v interface{})) {
	for !s.IsEmpty() {
		v := s.Get()
		for _, procedure := range procedures {
			procedure(v)
		}
	}
}

// Equals compares stack 's' with stack 'other' and returns true if they are equals.
// Time complexity: O(n), where n is the current length of 's' as well as 'other'.
func (s *Stack) Equals(other *Stack) bool {
	if s.len != other.len {
		return false
	}
	for i, j := s.top, other.top; i != nil; i, j = i.next, j.next {
		if i.value != j.value {
			return false
		}
	}
	return true
}

// EqualsByComparator compares stack 's' with stack 'other' and returns true if they are equals.
// The comparison between values is defined by parameter 'equals'.
// 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of 's' as well as 'other'.
func (s *Stack) EqualsByComparator(other *Stack, equals func(v1, v2 interface{}) bool) bool {
	if s.len != other.len {
		return false
	}
	for i, j := s.top, other.top; i != nil; i, j = i.next, j.next {
		if !equals(i.value, j.value) {
			return false
		}
	}
	return true
}

// Get returns the top value and removes it from the stack.
// If the stack is empty returns nil.
// Time complexity: O(1).
func (s *Stack) Get() interface{} {
	if s.IsEmpty() {
		return nil
	}
	n, v := s.top, s.top.value
	s.top = n.next
	n.clear()
	s.len--
	return v
}

// GetIf returns all first values that meet the condition defined by the "condition" parameter. These values
//will be removed from the queue.
// If the stack is empty returns an empty slice.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) GetIf(condition func(v interface{}) bool) []interface{} {
	values := make([]interface{}, 0)
	for n := s.top; n != nil; {
		if condition(n.value) {
			values = append(values, s.Get())
			n = s.top
		} else {
			return values
		}
	}
	return values
}

// IsEmpty returns true if the stack has no elements.
// Time complexity: O(1).
func (s *Stack) IsEmpty() bool {
	return s.len == 0
}

// Len returns the current length.
// Time complexity: O(1).
func (s *Stack) Len() int {
	return s.len
}

// Peek returns the top value.
// If the stack is empty returns nil.
// Time complexity: O(1).
func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return s.top.value
}

// Push inserts the value 'v' in the top.
// Time complexity: O(1).
func (s *Stack) Push(v interface{}) {
	n := &node{value: v, next: s.top}
	s.top = n
	s.len++
}

// RemoveAll sets the properties of this stack to its zero values.
// Time complexity: O(1).
func (s *Stack) RemoveAll() {
	s.top, s.len = nil, 0
}

// Search returns the index (zero based with top == s.Len()-1) of the first match of the value 'v'.
// If the value 'v' does not belong to the stack, return -1.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) Search(v interface{}) int {
	for n, i := s.top, s.len-1; n != nil; n, i = n.next, i-1 {
		if n.value == v {
			return i
		}
	}
	return -1
}

// SearchByComparator returns the index (zero based with top == s.Len()-1) of the first match of the value 'v'.
// If the value 'v' does not belong to the stack, return -1.
// The comparison between values is defined by parameter 'equals'.
// 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) SearchByComparator(v interface{}, equals func(v1, v2 interface{}) bool) int {
	for n, i := s.top, s.len-1; n != nil; n, i = n.next, i-1 {
		if equals(n.value, v) {
			return i
		}
	}
	return -1
}

// Slice returns a new slice with the values stored in the stack keeping its order.
// 's' retains its original state.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) Slice() []interface{} {
	values := make([]interface{}, 0, s.len)
	for n := s.top; n != nil; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// String returns a representation of the stack as a string.
// Stack implements the fmt.Stringer interface.
// Time complexity: O(n), where n is the current length of 's'.
func (s *Stack) String() string {
	if s.IsEmpty() {
		return "[]"
	}
	str := "["
	n := s.top
	for ; n.next != nil; n = n.next {
		str += fmt.Sprintf("%v ", n.value)
	}
	return str + fmt.Sprintf("%v]", n.value)
}
