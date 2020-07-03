// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package sortedset implements an AVL tree with ordered set behaviors.
package sortedset

import "fmt"

// node of a AVL tree.
type node struct {
	// value stored in the node.
	value interface{}

	// left points to a node smaller than this.
	// right points to a node greater than this.
	left, right *node

	// h is the current height (number of levels in the AVL tree).
	// len is the current length (number of nodes).
	h, len int
}

// clear sets the properties of the node to its zero values.
// Time complexity: O(1).
func (n *node) clear() {
	n.value, n.left, n.right, n.h, n.len = nil, nil, nil, 0, 0
}

// SortedSet represents a AVL tree.
// The zero value for SortedSet is an empty SortedSet ready to use.
type SortedSet struct {
	// root points to the root node in the AVL tree.
	root *node
}

// balance returns the balance of the particular AVL tree 'n'.
// If 'n' equals nil, then return 0.
// Time complexity: O(1).
func balance(n *node) int {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}

// height returns the height of the particular AVL tree 'n'.
// If 'n' equals nil, then return 0.
// Time complexity: O(1).
func height(n *node) int {
	if n == nil {
		return 0
	}
	return n.h
}

// leftRotate do the avl tree left rotate with 'n' as a root.
// Time complexity: O(1).
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

// length returns the length of the particular AVL tree 'n'.
// If 'n' equals nil, then return 0.
// Time complexity: O(1).
func length(n *node) int {
	if n == nil {
		return 0
	}
	return n.len
}

// maxInt returns the biggest integer between 'a' and 'b'.
// Time complexity: O(1).
func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// min returns the minimum value in the avl tree 'n'.
// Time complexity: O(log(n)), where n is the current length of the AVL tree.
func min(n *node) *node {
	if n == nil {
		return nil
	}
	for n.left != nil {
		n = n.left
	}
	return n
}

// rightRotate do the avl tree right rotate with 'n' as a root.
// Time complexity: O(1).
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
// Time complexity: O(1).
func New() *SortedSet {
	return new(SortedSet)
}

// NewBySlice returns a new SortedSet with the values stored in the slice.
// The comparison to order the values is defined by the parameter 'compare'.
// The function 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or
//greater than 'v2'.
// Time complexity: O(n*log(n)), where n is the current length of the slice.
func NewBySlice(values []interface{}, compare func(v1, v2 interface{}) int) *SortedSet {
	s := New()
	for _, v := range values {
		s.Push(v, compare)
	}
	return s
}

// Clone returns a new cloned SortedSet.
// Time complexity: O(n), where n is the current length of the set.
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

// Contains returns true if the value 'v' belongs to the set.
// The comparison to order the values is defined by the parameter 'compare'.
// The function 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or
//greater than 'v2'.
// Time complexity: O(log(n)), where n is the current length of the set.
func (s *SortedSet) Contains(v interface{}, compare func(v1, v2 interface{}) int) bool {
	return containsRecursive(v, s.root, compare)
}

// containsRecursive is an auxiliary recursive function of the SortedSet Contains method.
func containsRecursive(v interface{}, n *node, compare func(v1, v2 interface{}) int) bool {
	if n == nil {
		return false
	}
	switch diff := compare(v, n.value); {
	case diff < 0:
		return containsRecursive(v, n.left, compare)
	case diff > 0:
		return containsRecursive(v, n.right, compare)
	default: // diff == 0
		return true
	}
}

// Do gets the first (minor) value and performs all the procedures, then repeats it with the rest of the values.
// The set retains its original state.
// Time complexity: O(n*p), where n is the current length of the set and p is the number of procedures.
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

// IsEmpty returns true if the set has no values.
// Time complexity: O(1).
func (s *SortedSet) IsEmpty() bool {
	return s.root == nil
}

// Len returns the current length of the set.
// Time complexity: O(1).
func (s *SortedSet) Len() int {
	if s.IsEmpty() {
		return 0
	}
	return s.root.len
}

// Max returns the maximum value of the set.
// Time complexity: O(log(n)), where n is the current length of the set.
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

// Min returns the minimum value of the set.
// Time complexity: O(log(n)), where n is the current length of the set.
func (s *SortedSet) Min() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return min(s.root).value
}

// Push inserts the value 'v' in an orderly way.
// The comparison to order the values is defined by the parameter 'compare'.
// The function 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or
//greater than 'v2'.
// Time complexity: O(log(n)), where n is the current length of the set.
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
// The comparison to order the values is defined by the parameter 'compare'.
// The function 'compare' must return a negative int, zero, or a positive int as 'v1' is less than, equal to, or
//greater than 'v2'.
// Time complexity: O(log(n)), where n is the current length of the set.
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

// RemoveAll sets the properties of the set to its zero values.
// Time complexity: O(1).
func (s *SortedSet) RemoveAll() {
	s.root = nil
}

// Slice returns a new slice with the values stored in the set keeping its order.
// The set retains its original state.
// Time complexity: O(n), where n is the current length of the set.
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
// Time complexity: O(n), where n is the current length of the set.
func (s *SortedSet) String() string {
	if s.IsEmpty() {
		return "[]"
	}
	str := "[" + stringRecursive(s.root)
	return str[:len(str)-1] + "]"
}

// stringRecursive is an auxiliary recursive function of the SortedSet String method.
func stringRecursive(n *node) string {
	if n == nil {
		return ""
	}
	return stringRecursive(n.left) + fmt.Sprintf("%v ", n.value) + stringRecursive(n.right)
}
