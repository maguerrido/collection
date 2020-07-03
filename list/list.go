// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package list implements a doubly-linked list.
package list

import "fmt"

// Element of a list.
type Element struct {
	// value stored in the element.
	value interface{}

	// next points to the next node.
	// prev points to the previous node.
	// If the element is the front or back element, then next or prev points to nil respectively.
	next, prev *Element

	// parent points to the list containing the element.
	parent *List
}

// clear sets the properties of the element to its zero values.
// Time complexity: O(1).
func (e *Element) clear() {
	e.value, e.next, e.prev, e.parent = nil, nil, nil, nil
}

// Next returns the next list element.
// If this node is the back node, then returns nil.
// Time complexity: O(1).
func (e *Element) Next() *Element {
	return e.next
}

// Parent returns the list containing this element.
// Time complexity: O(1).
func (e *Element) Parent() *List {
	return e.parent
}

// Prev returns the previous list element.
// If this node is the front node, then returns nil.
// Time complexity: O(1).
func (e *Element) Prev() *Element {
	return e.prev
}

// Set updates the value stored in this element.
// Time complexity: O(1).
func (e *Element) Set(v interface{}) {
	e.value = v
}

// Value returns the value stored in this element.
// Time complexity: O(1).
func (e *Element) Value() interface{} {
	return e.value
}

// List represents a doubly-linked list.
// The zero value for List is an empty List ready to use.
type List struct {
	// front points to the front (first) element in the list.
	// back points to the back (last) element in the list.
	front, back *Element

	// len is the current length (number of elements).
	len int
}

// New returns a new List ready to use.
// Time complexity: O(1).
func New() *List {
	return new(List)
}

// NewBySlice returns a new List with the values stored in the slice keeping its order.
// Time complexity: O(n), where n is the current length of the slice.
func NewBySlice(values []interface{}) *List {
	l := New()
	for _, v := range values {
		l.PushBack(v)
	}
	return l
}

// Back returns the back element.
// If the list is empty, then returns nil.
// Time complexity: O(1).
func (l *List) Back() *Element {
	if l.IsEmpty() {
		return nil
	}
	return l.back
}

// Clone returns a new cloned List.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) Clone() *List {
	clone := New()
	for e := l.front; e != nil; e = e.next {
		clone.PushBack(e.value)
	}
	return clone
}

// Contains returns true if the element 'e' belongs to the list.
// Time complexity: O(1).
func (l *List) Contains(e *Element) bool {
	return e != nil && e.parent == l
}

// Do gets the front value and performs all the procedures, then repeats it with the rest of the values.
// The list retains its original state.
// Time complexity: O(n*p), where n is the current length of the list and p is the number of procedures.
func (l *List) Do(procedures ...func(v interface{})) {
	for e := l.front; e != nil; e = e.next {
		for _, procedure := range procedures {
			procedure(e.value)
		}
	}
}

// Equals compares this list with the 'other' list and returns true if they are equal.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) Equals(other *List) bool {
	if l.len != other.len {
		return false
	}
	for i, j := l.front, other.front; i != nil; i, j = i.next, j.next {
		if i.value != j.value {
			return false
		}
	}
	return true
}

// EqualsByComparator compares this list with the 'other' list and returns true if they are equal.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) EqualsByComparator(other *List, equals func(v1, v2 interface{}) bool) bool {
	if l.len != other.len {
		return false
	}
	for i, j := l.front, other.front; i != nil; i, j = i.next, j.next {
		if !equals(i.value, j.value) {
			return false
		}
	}
	return true
}

// Front returns the front element.
// If the list is empty, then returns nil.
// Time complexity: O(1).
func (l *List) Front() *Element {
	if l.IsEmpty() {
		return nil
	}
	return l.front
}

// Get returns the 'index' (zero based) position element.
// If 'index' is out of bounds, then returns nil.
// Time complexity: O(n/2), where n is the current length of the list.
func (l *List) Get(index int) *Element {
	if index < 0 || index > l.len-1 {
		return nil
	}
	var e *Element
	if index < l.len/2 {
		e = l.front
		for i := 0; i < index; i++ {
			e = e.next
		}
	} else {
		e = l.back
		for i := l.len - 1; i > index; i-- {
			e = e.prev
		}
	}
	return e
}

// IsEmpty returns true if the list has no elements.
// Time complexity: O(1).
func (l *List) IsEmpty() bool {
	return l.len == 0
}

// Len returns the current length of the list.
// Time complexity: O(1).
func (l *List) Len() int {
	return l.len
}

// MoveAfter moves the element 'e' after 'mark'.
// Time complexity: O(1).
func (l *List) MoveAfter(e, mark *Element) bool {
	if !l.Contains(e) || !l.Contains(mark) {
		return false
	}
	if e == mark {
		return true
	}
	l.unlink(e)
	if mark.next == nil {
		l.back = e
	} else {
		mark.next.prev = e
	}
	e.next = mark.next
	e.prev = mark
	mark.next = e
	return true
}

// MoveBefore moves the element 'e' before 'mark'.
// Time complexity: O(1).
func (l *List) MoveBefore(e, mark *Element) bool {
	if !l.Contains(e) || !l.Contains(mark) {
		return false
	}
	if e == mark {
		return true
	}
	l.unlink(e)
	if mark.prev == nil {
		l.front = e
	} else {
		mark.prev.next = e
	}
	e.prev = mark.prev
	e.next = mark
	mark.prev = e
	return true
}

// MoveToBack moves the element 'e' to back.
// Time complexity: O(1).
func (l *List) MoveToBack(e *Element) bool {
	return l.MoveAfter(e, l.back)
}

// MoveToFront moves the element 'e' to front.
// Time complexity: O(1).
func (l *List) MoveToFront(e *Element) bool {
	return l.MoveAfter(e, l.front)
}

// PushAfter inserts the value 'v' after the element 'mark'.
// Time complexity: O(1).
func (l *List) PushAfter(v interface{}, mark *Element) *Element {
	if !l.Contains(mark) {
		return nil
	}
	e := &Element{value: v, next: mark.next, prev: mark, parent: l}
	mark.next = e
	if mark == l.back {
		l.back = e
	} else {
		e.next.prev = e
	}
	l.len++
	return e
}

// PushBack inserts the value 'v' at the back of the list.
// Time complexity: O(1).
func (l *List) PushBack(v interface{}) {
	e := &Element{value: v, next: nil, prev: l.back, parent: l}
	if l.IsEmpty() {
		l.front = e
	} else {
		l.back.next = e
	}
	l.back = e
	l.len++
}

// PushBackList inserts the list 'other' at the back of this list.
// Time complexity: O(n), where n is the current length of the list 'other'.
func (l *List) PushBackList(other *List) {
	if other == nil {
		return
	}
	for e := other.front; e != nil; e = e.next {
		l.PushBack(e.value)
	}
}

// PushBefore inserts the value 'v' before the element 'mark'.
// Time complexity: O(1).
func (l *List) PushBefore(v interface{}, mark *Element) *Element {
	if !l.Contains(mark) {
		return nil
	}
	e := &Element{value: v, next: mark, prev: mark.prev, parent: l}
	mark.prev = e
	if mark == l.front {
		l.front = e
	} else {
		e.prev.next = e
	}
	l.len++
	return e
}

// PushFront inserts the value 'v' at the front of the list.
// Time complexity: O(1).
func (l *List) PushFront(v interface{}) {
	e := &Element{value: v, next: l.front, prev: nil, parent: l}
	if l.IsEmpty() {
		l.back = e
	} else {
		l.front.prev = e
	}
	l.front = e
	l.len++
}

// PushFrontList inserts the list 'other' in the front of this list.
// Time complexity: O(n), where n is the current length of the list 'other'.
func (l *List) PushFrontList(other *List) {
	if other == nil {
		return
	}
	for e := other.back; e != nil; e = e.prev {
		l.PushFront(e.value)
	}
}

// quickSort sorts the list using the Quick Sort algorithm.
func (l *List) quickSort(compare func(v1, v2 interface{}) int) {
	quickSortRecursive(l.front, l.back, compare)
}

// quickSortRecursive is an auxiliary recursive function of the List quickSort method.
func quickSortRecursive(front, back *Element, compare func(v1, v2 interface{}) int) {
	if front != nil && back != nil && front != back.next {
		pivot := front.prev
		for j := front; j != back; j = j.next {
			if compare(j.value, back.value) < 1 {
				if pivot == nil {
					pivot = front
				} else {
					pivot = pivot.next
				}
				pivot.value, j.value = j.value, pivot.value // swap
			}
		}
		if pivot == nil {
			pivot = front
		} else {
			pivot = pivot.next
		}
		pivot.value, back.value = back.value, pivot.value // swap
		quickSortRecursive(front, pivot.prev, compare)
		quickSortRecursive(pivot.next, back, compare)
	}
}

// Remove removes the first match of the value 'v' in the list.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) Remove(v interface{}) bool {
	for e := l.front; e != nil; e = e.next {
		if v == e.value {
			l.unlink(e)
			e.clear()
			l.len--
			return true
		}
	}
	return false
}

// RemoveAll sets the properties of the list to its zero values.
// Time complexity: O(1).
func (l *List) RemoveAll() {
	l.front, l.back, l.len = nil, nil, 0
}

// RemoveElement removes the element 'e' from the list.
// Time complexity: O(1).
func (l *List) RemoveElement(e *Element) (v interface{}, ok bool) {
	if !l.Contains(e) {
		return nil, false
	}
	l.unlink(e)
	l.len--
	v = e.value
	e.clear()
	return v, true
}

// RemoveFrom removes the element 'start', 'end' and all elements between them, and then returns the number of removals.
// If the element 'end' is before 'start', then 'start' and all the next elements will be removed.
// Time complexity: O(n), where n is the length between 'start' and 'end'.
func (l *List) RemoveFrom(start, end *Element) int {
	if !l.Contains(start) || !l.Contains(end) {
		return 0
	}
	count := 0
	for start != end && start != nil {
		toRemove := start
		start = start.next
		l.RemoveElement(toRemove)
		count++
	}
	if start == end {
		l.RemoveElement(start)
		count++
	}
	return count
}

// RemoveIf removes all values that meet the condition defined by the parameter 'condition' and returns the number of
//removals.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) RemoveIf(condition func(v interface{}) bool) int {
	count := 0
	for e := l.front; e != nil; {
		if condition(e.value) {
			next := e.next
			l.unlink(e)
			e.clear()
			l.len--
			count++
			e = next
		} else {
			e = e.next
		}
	}
	return count
}

// Search returns the index (zero based) of the first match of the value 'v' and the element containing it.
// If the value 'v' does not belong to the list, then returns -1 and nil.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) Search(v interface{}) (index int, e *Element) {
	for e, i := l.front, 0; e != nil; e, i = e.next, i+1 {
		if e.value == v {
			return i, e
		}
	}
	return -1, nil
}

// SearchByComparator returns the index (zero based) of the first match of the value 'v' and the element containing it.
// If the value 'v' does not belong to the list, then returns -1 and nil.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) SearchByComparator(v interface{}, equals func(v1, v2 interface{}) bool) (index int, e *Element) {
	for e, i := l.front, 0; e != nil; e, i = e.next, i+1 {
		if equals(e.value, v) {
			return i, e
		}
	}
	return -1, nil
}

// selectionSort sorts the list using the Selection Sort algorithm.
func (l *List) selectionSort(compare func(v1, v2 interface{}) int) {
	for i := l.front; i != l.back; i = i.next {
		minor := i
		for j := minor.next; j != nil; j = j.next {
			if compare(j.value, minor.value) < 1 {
				minor = j
			}
		}
		i.value, minor.value = minor.value, i.value // swap
	}
}

// Slice returns a new slice with the values stored in the list keeping its order.
// The list retains its original state.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) Slice() []interface{} {
	values := make([]interface{}, 0, l.len)
	for e := l.front; e != nil; e = e.next {
		values = append(values, e.value)
	}
	return values
}

// Sort sorts the list.
// The comparison to order the values is defined by the parameter 'compare'.
// The function 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or
//greater than 'v2'.
// Time complexity (n <= 10): O(n^2), where n is the current length of the list.
// Time complexity (n > 10): θ(n*log(n)) and O(n^2), where n is the current length of the list.
func (l *List) Sort(compare func(v1, v2 interface{}) int) {
	if l.len > 1 {
		if l.len > 10 {
			l.quickSort(compare)
		} else {
			l.selectionSort(compare)
		}
	}
}

// String returns a representation of the list as a string.
// List implements the fmt.Stringer interface.
// Time complexity: O(n), where n is the current length of the list.
func (l *List) String() string {
	if l.IsEmpty() {
		return "[]"
	}
	str := "["
	e := l.front
	for ; e.next != nil; e = e.next {
		str += fmt.Sprintf("%v ", e.value)
	}
	return str + fmt.Sprintf("%v]", e.value)
}

// Swap swaps the values between elements 'a' and 'b'.
// Time complexity: O(1).
func (l *List) Swap(a, b *Element) bool {
	if !l.Contains(a) || !l.Contains(b) {
		return false
	}
	a.value, b.value = b.value, a.value
	return true
}

// unlink unlinks an element in the list.
// Time complexity: O(1).
func (l *List) unlink(e *Element) {
	if e.prev == nil {
		if e.next == nil {
			l.front, l.back = nil, nil
		} else { // e.next != nil
			e.next.prev = nil
			l.front = e.next
		}
	} else { // e.prev != nil
		if e.next != nil {
			e.prev.next = e.next
			e.next.prev = e.prev
		} else { // e.next == nil
			e.prev.next = nil
			l.back = e.prev
		}
	}
	e.next = nil
	e.prev = nil
}
