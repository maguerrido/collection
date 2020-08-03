// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package stack implements a singly-linked list with stack behaviors.
package stack

import (
	"fmt"
	coll "github.com/maguerrido/collection"
)

// node of a Stack.
type node struct {
	// value stored in the node.
	value interface{}

	// next points to the next node.
	// If the node is the bottom node, then points to nil.
	next *node
}

// clear sets the properties of the node to its zero values.
// Time complexity: O(1).
func (n *node) clear() {
	n.value, n.next = nil, nil
}

// Stack represents a singly-linked list.
// The zero value for Stack is an empty Stack ready to use.
type Stack struct {
	// top points to the top node in the Stack.
	top *node

	// len is the current length (number of nodes).
	len int
}

// New returns a new Stack ready to use.
// Time complexity: O(1).
func New() *Stack {
	return new(Stack)
}

// NewBySlice returns a new Stack with the values stored in the slice keeping its order.
// The last value of the slice will be the top value of the stack.
// Time complexity: O(n), where n is the current length of the slice.
func NewBySlice(values []interface{}) *Stack {
	s := New()
	for _, v := range values {
		s.Push(v)
	}
	return s
}

// Clone returns a new cloned Stack.
// Time complexity: O(n), where n is the current length of the stack.
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

// Do gets the top value and performs all the procedures, then repeats it with the rest of the values.
// The stack will be empty.
// Time complexity: O(n*p), where n is the current length of the stack and p is the number of procedures.
func (s *Stack) Do(procedures ...func(v interface{})) {
	for !s.IsEmpty() {
		v := s.Get()
		for _, procedure := range procedures {
			procedure(v)
		}
	}
}

// Equals compares this stack with the 'other' stack and returns true if they are equal.
// Time complexity: O(n), where n is the current length of the stack.
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

// EqualsByComparator compares this stack with the 'other' stack and returns true if they are equal.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of the stack.
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
// If the stack is empty, then returns nil.
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

// GetIf returns all first values that meet the condition defined by the 'condition' parameter. These values will be
//removed from the stack.
// Time complexity: O(n), where n is the current length of the stack.
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

// IsEmpty returns true if the stack has no values.
// Time complexity: O(1).
func (s *Stack) IsEmpty() bool {
	return s.len == 0
}

func (s *Stack) Iterator() coll.Iterator {
	return &iterator{
		s:           s,
		prev:        nil,
		this:        nil,
		index:       -1,
		lastCommand: -1,
		lastHasNext: false,
	}
}

// Len returns the current length of the stack.
// Time complexity: O(1).
func (s *Stack) Len() int {
	return s.len
}

// Peek returns the top value.
// If the stack is empty, then returns nil.
// Time complexity: O(1).
func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return s.top.value
}

// Push inserts the value 'v' at the top of the stack.
// Time complexity: O(1).
func (s *Stack) Push(v interface{}) {
	n := &node{value: v, next: s.top}
	s.top = n
	s.len++
}

// RemoveAll sets the properties of the stack to its zero values.
// Time complexity: O(1).
func (s *Stack) RemoveAll() {
	s.top, s.len = nil, 0
}

func (s *Stack) removeNode(prev, n *node) {
	if s.top == n {
		s.top = n.next
	} else {
		prev.next = n.next
	}
	n.clear()
	s.len--
}

// Search returns the index (zero based with top equal to current length - 1) of the first match of the value 'v'.
// If the value 'v' does not belong to the stack, then returns -1.
// Time complexity: O(n), where n is the current length of the stack.
func (s *Stack) Search(v interface{}) int {
	for n, i := s.top, s.len-1; n != nil; n, i = n.next, i-1 {
		if n.value == v {
			return i
		}
	}
	return -1
}

// SearchByComparator returns the index (zero based with top equal to current length - 1) of the first match of the
//value 'v'.
// If the value 'v' does not belong to the stack, then returns -1.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of the stack.
func (s *Stack) SearchByComparator(v interface{}, equals func(v1, v2 interface{}) bool) int {
	for n, i := s.top, s.len-1; n != nil; n, i = n.next, i-1 {
		if equals(n.value, v) {
			return i
		}
	}
	return -1
}

// Slice returns a new slice with the values stored in the stack from top to bottom.
// The stack retains its original state.
// Time complexity: O(n), where n is the current length of the stack.
func (s *Stack) Slice() []interface{} {
	values := make([]interface{}, 0, s.len)
	for n := s.top; n != nil; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// String returns a representation of the stack as a string.
// Stack implements the fmt.Stringer interface.
// Time complexity: O(n), where n is the current length of the stack.
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

type iterator struct {
	s           *Stack
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
		for n := i.s.top; n != nil; n = n.next {
			action(&n.value)
		}
	}
}

func (i *iterator) HasNext() bool {
	i.lastCommand = iteratorCommandHasNext
	i.lastHasNext = i.index < i.s.len-1
	return i.lastHasNext
}

func (i *iterator) Next() (interface{}, error) {
	if i.lastCommand != iteratorCommandHasNext {
		return nil, fmt.Errorf(coll.ErrorIteratorNext)
	} else if !i.lastHasNext {
		return nil, fmt.Errorf(coll.ErrorIteratorHasNext)
	}

	if i.this == nil {
		i.this = i.s.top
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

	i.s.removeNode(i.prev, i.this)
	i.this = i.prev
	i.index--
	i.lastCommand = iteratorCommandRemove

	return nil
}
