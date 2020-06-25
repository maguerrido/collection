// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package sortedset implements a ordered set using an AVL tree.
package sortedset

import (
	"fmt"
	"strings"
)

// node represents a binary tree.
type node struct {
	// value stored in this node.
	value interface{}

	// left points to a node smaller than this.
	// right points to a node greater than this.
	left, right *node

	// h is the current height.
	// len is the current height.
	h, len int
}

// clear sets the properties of this node to its zero values.
func (n *node) clear() {
	n.value, n.left, n.right, n.h, n.len = nil, nil, nil, 0, 0
}

// SortedSet represents a AVL tree.
// The zero value for SortedSet is an empty SortedSet ready to use.
type SortedSet struct {
	// root points to the root node in the SortedSet.
	root *node
}

func balance(n *node) int {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}
func height(n *node) int {
	if n == nil {
		return 0
	}
	return n.h
}
func leftRotate(n *node) *node {
	root := n.right
	tree := root.left
	root.left = n
	n.right = tree
	n.h = 1 + maxInt(height(n.left), height(n.right))
	n.len = 1 + length(n.left) + length(n.right)
	root.h = 1 + maxInt(height(root.left), height(root.right))
	root.len = 1 + length(root.left) + length(root.right)
	return root
}
func length(n *node) int {
	if n == nil {
		return 0
	}
	return n.len
}
func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
func min(n *node) *node {
	if n == nil {
		return nil
	}
	for n.left != nil {
		n = n.left
	}
	return n
}
func rightRotate(n *node) *node {
	root := n.left
	tree := root.right
	root.right = n
	n.left = tree
	n.h = 1 + maxInt(height(n.left), height(n.right))
	n.len = 1 + length(n.left) + length(n.right)
	root.h = 1 + maxInt(height(root.left), height(root.right))
	root.len = 1 + length(root.left) + length(root.right)
	return root
}

// New returns a new SortedSet ready to use.
func New() *SortedSet {
	return new(SortedSet)
}

// NewBySlice returns a new SortedSet with the values stored in the slice.
// The comparison between values is defined by parameter 'compare'.
// 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or greater than 'v2'.
// Time complexity: O(n * log n), where n is the current length of 'values'.
func NewBySlice(values []interface{}, compare func(v1, v2 interface{}) int) *SortedSet {
	q := New()
	for _, v := range values {
		q.Push(v, compare)
	}
	return q
}

// Clone returns a new cloned SortedSet of 's'.
// 's' retains its original state.
// Time complexity: O(n), where n is the current length of 's'.
func (s *SortedSet) Clone() *SortedSet {
	clone := New()
	clone.root = cloneRecursive(s.root)
	return clone
}

// cloneRecursive is an auxiliary recursive function of the SortedSet Clone method.
func cloneRecursive(nS *node) *node {
	if nS == nil {
		return nil
	}
	left := cloneRecursive(nS.left)
	right := cloneRecursive(nS.right)
	return &node{value: nS.value, left: left, right: right, h: nS.h, len: nS.len}
}

// Contains returns true if 'v' belongs to 's'.
// Time complexity: O(n), where n is the current length of 's'.
func (s *SortedSet) Contains(v interface{}) bool {
	return containsRecursive(v, s.root)
}

// containsRecursive is an auxiliary recursive function of the SortedSet Contains method.
func containsRecursive(v interface{}, n *node) bool {
	if n == nil {
		return false
	} else if n.value == v {
		return true
	} else if containsRecursive(v, n.left) {
		return true
	} else {
		return containsRecursive(v, n.right)
	}
}

// Contains returns true if 'v' belongs to 's'.
// The comparison between values is defined by parameter 'compare'.
// 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or greater than 'v2'.
// Time complexity: O(log n), where n is the current length of 's'.
func (s *SortedSet) ContainsByComparator(v interface{}, compare func(v1, v2 interface{}) int) bool {
	return containsByComparatorRecursive(v, s.root, compare)
}

// containsByComparatorRecursive is an auxiliary recursive function of the SortedSet ContainsByComparator method.
func containsByComparatorRecursive(v interface{}, n *node, compare func(v1, v2 interface{}) int) bool {
	if n == nil {
		return false
	}
	switch diff := compare(v, n.value); {
	case diff < 0:
		return containsByComparatorRecursive(v, n.left, compare)
	case diff > 0:
		return containsByComparatorRecursive(v, n.right, compare)
	default: // diff == 0
		return true
	}
}

// Do gets the first (minor) value and performs all the procedures, then repeats this with the rest of the values.
// Time complexity: O(n), where n is the current length of 's'.
func (s *SortedSet) Do(procedures ...func(v interface{})) {
	doRecursive(s.root, procedures...)
}

// doRecursive is an auxiliary recursive function of the SortedSet Do method.
func doRecursive(n *node, procedures ...func(v interface{})) {
	if n == nil {
		return
	}
	doRecursive(n.left, procedures...)
	for _, procedure := range procedures {
		procedure(n.value)
	}
	doRecursive(n.right, procedures...)
}

// IsEmpty returns true if the set has no elements.
// Time complexity: O(1).
func (s *SortedSet) IsEmpty() bool {
	return s.root == nil
}

// Len returns the current length.
// Time complexity: O(1).
func (s *SortedSet) Len() int {
	if s.IsEmpty() {
		return 0
	}
	return s.root.len
}

// Max returns the maximum value.
// Time complexity: O(log n), where n is the current length of 's'.
func (s *SortedSet) Max() interface{} {
	if s.IsEmpty() {
		return nil
	}
	n := s.root
	for n.right != nil {
		n = n.right
	}
	return n.value
}

// Min returns the minimum value.
// Time complexity: O(log n), where n is the current length of 's'.
func (s *SortedSet) Min() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return min(s.root).value
}

// Push inserts the value 'v' in order.
// If 'v' already exists, then it will be updated.
// The comparison between values is defined by parameter 'compare'.
// 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or greater than 'v2'.
// Time complexity: O(log n), where n is the current length of 's'.
func (s *SortedSet) Push(v interface{}, compare func(v1, v2 interface{}) int) {
	s.root = pushRecursive(v, s.root, compare)
}

// pushRecursive is an auxiliary recursive function of the SortedSet Push method.
func pushRecursive(v interface{}, n *node, compare func(v1, v2 interface{}) int) *node {
	if n == nil {
		return &node{value: v, left: nil, right: nil, h: 1, len: 1}
	}

	switch diff := compare(v, n.value); {
	case diff < 0:
		n.left = pushRecursive(v, n.left, compare)
	case diff > 0:
		n.right = pushRecursive(v, n.right, compare)
	case diff == 0:
		n.value = v
		return n
	}

	n.h = 1 + maxInt(height(n.left), height(n.right))
	n.len = 1 + length(n.left) + length(n.right)

	balance := balance(n)
	if balance > 1 {
		if compare(v, n.left.value) < 0 { // case: left-left
			return rightRotate(n)
		} else { // case: left-right
			n.left = leftRotate(n.left)
			return rightRotate(n)
		}
	}
	if balance < -1 {
		if compare(v, n.right.value) > 0 { // case: right-right
			return leftRotate(n)
		} else { // case: right-left
			n.right = rightRotate(n.right)
			return leftRotate(n)
		}
	}

	return n
}

// Remove removes the value 'v' from the set.
// The comparison between values is defined by parameter 'compare'.
// 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or greater than 'v2'.
// Time complexity: O(log n), where n is the current length of 's'.
func (s *SortedSet) Remove(v interface{}, compare func(v1, v2 interface{}) int) bool {
	var removed bool
	s.root, removed = removeRecursive(v, s.root, compare)
	return removed
}

// removeRecursive is an auxiliary recursive function of the SortedSet Remove method.
func removeRecursive(v interface{}, n *node, compare func(v1, v2 interface{}) int) (*node, bool) {
	if n == nil {
		return nil, false
	}

	removed := false
	switch diff := compare(v, n.value); {
	case diff < 0:
		n.left, removed = removeRecursive(v, n.left, compare)
	case diff > 0:
		n.right, removed = removeRecursive(v, n.right, compare)
	case diff == 0:
		removed = true
		var temp *node
		// case: node with only one child or no child
		if n.left == nil || n.right == nil {
			if n.left == nil {
				temp = n.right
			} else {
				temp = n.left
			}
			if temp == nil { // case: no child
				temp = n
				n.clear()
				n = nil
			} else { // case: one child
				n.value = temp.value
				temp.clear()
			}
		} else { // case: node with two children
			temp = min(n.right)
			n.value = temp.value
			n.right, _ = removeRecursive(temp.value, n.right, compare)
		}
	}

	// case: no child (removed)
	if n == nil {
		return nil, removed
	}

	n.h = 1 + maxInt(height(n.left), height(n.right))
	n.len = 1 + length(n.left) + length(n.right)

	balanceTreeNode := balance(n)
	if balanceTreeNode > 1 {
		if balance(n.left) >= 0 { // case: left-left
			return rightRotate(n), removed
		} else { // case: left-right
			n.left = leftRotate(n.left)
			return rightRotate(n), removed
		}
	}
	if balanceTreeNode < -1 {
		if balance(n.right) <= 0 { // case: right-right
			return leftRotate(n), removed
		} else { // case: right-left
			n.right = rightRotate(n.right)
			return leftRotate(n), removed
		}
	}

	return n, removed
}

// RemoveAll sets the properties of this set to its zero values.
// Time complexity: O(1).
func (s *SortedSet) RemoveAll() {
	s.root = nil
}

// Slice returns a new slice with the values stored in the set keeping its order.
// 's' retains its original state.
// Time complexity: O(n), where n is the current length of 's'.
func (s *SortedSet) Slice() []interface{} {
	values := make([]interface{}, 0, s.Len())
	sliceRecursive(s.root, &values)
	return values
}

// sliceRecursive is an auxiliary recursive function of the SortedSet Slice method.
func sliceRecursive(n *node, values *[]interface{}) {
	if n == nil {
		return
	}
	sliceRecursive(n.left, values)
	*values = append(*values, n.value)
	sliceRecursive(n.right, values)
}

// String returns a representation of the set as a string.
// SortedSet implements the fmt.Stringer interface.
// Time complexity: O(n), where n is the current length of 's'.
func (s *SortedSet) String() string {
	str := "["
	if s.IsEmpty() {
		str += "]"
	} else {
		str += stringRecursive(s.root)
		str = strings.TrimRight(str, " ")
		str += "]"
	}
	return str
}

// stringRecursive is an auxiliary recursive function of the SortedSet String method.
func stringRecursive(n *node) string {
	if n == nil {
		return ""
	}
	return stringRecursive(n.left) + fmt.Sprintf("%v ", n.value) + stringRecursive(n.right)
}
