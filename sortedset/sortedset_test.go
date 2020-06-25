// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package sortedset

import (
	"fmt"
	"testing"
)

func checkHeight(n *node) bool {
	if n == nil {
		return true
	}
	checkHeight(n.left)
	checkHeight(n.right)
	if n.h != 1+maxInt(heightRecursive(n.left), heightRecursive(n.right)) {
		return false
	}
	return true
}
func checkLength(n *node) bool {
	if n == nil {
		return true
	}
	checkLength(n.left)
	checkLength(n.right)
	if n.len != 1+lengthRecursive(n.left)+lengthRecursive(n.right) {
		return false
	}
	return true
}
func checkOrder(n *node, compare func(v1, v2 interface{}) int) bool {
	if n == nil {
		return true
	}
	bLeft, bRight := true, true
	if n.left != nil {
		bLeft = compare(n.value, n.left.value) > 0 && checkOrder(n.left, compare)
	}
	if n.right != nil {
		bRight = compare(n.value, n.right.value) < 0 && checkOrder(n.right, compare)
	}
	return bLeft && bRight
}
func checkValues(n *node, values []interface{}) bool {
	if n == nil {
		if len(values) != 0 {
			return false
		} else {
			return true
		}
	} else {
		if len(values) == 0 {
			return false
		}
	}
	for _, v := range values {
		if !containsRecursive(v, n) {
			return false
		}
	}
	return true
}
func checkValuesAndPositions(rootS, rootOther *node) bool {
	if rootS == nil {
		if rootOther == nil {
			return true
		} else {
			return false
		}
	}
	if rootOther == nil {
		return false
	}
	return rootS.value == rootOther.value && checkValuesAndPositions(rootS.left, rootOther.left) && checkValuesAndPositions(rootS.right, rootOther.right)
}
func checkZeroValue(s *SortedSet) bool {
	return s.root == nil
}
func compareInt(v1, v2 interface{}) int {
	int1 := v1.(int)
	int2 := v2.(int)
	return int1 - int2
}
func heightRecursive(n *node) int {
	if n == nil {
		return 0
	}
	return 1 + maxInt(heightRecursive(n.left), heightRecursive(n.right))
}
func lengthRecursive(n *node) int {
	if n == nil {
		return 0
	}
	return 1 + lengthRecursive(n.left) + lengthRecursive(n.right)
}
func sortedset(len int) *SortedSet {
	s := New()
	for i := 0; i < len; i++ {
		s.Push(i, compareInt)
	}
	return s
}

func TestNew(t *testing.T) {
	got := New()
	if !checkZeroValue(got) {
		t.Errorf("checkZeroValue: FAIL")
	}
}
func TestNewBySlice(t *testing.T) {
	tests := []struct {
		name string
		in   []interface{}
	}{
		{"empty", []interface{}{}},
		{"!empty", []interface{}{5, 2, 1, 10, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			s := NewBySlice(test.in, compareInt)
			if !checkValues(s.root, test.in) {
				tt.Errorf("checkValues: FAIL")
			}
		})
	}
}

func TestSortedSet_Clone(t *testing.T) {
	tests := []struct {
		name string
		s    *SortedSet
	}{
		{"empty", New()},
		{"!empty", NewBySlice([]interface{}{5, 2, 3, 6, 8, 1, 10, 4}, compareInt)},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			clone := test.s.Clone()
			if !checkValuesAndPositions(test.s.root, clone.root) {
				tt.Errorf("checkValuesAndPositions: FAIL")
			}
		})
	}
}
func TestSortedSet_Contains(t *testing.T) {
	tests := []struct {
		name string
		s    *SortedSet
		in   int
		out  bool
	}{
		{"empty", New(), 1, false},
		{"!empty/false", sortedset(10), 20, false},
		{"!empty/true", sortedset(10), 5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Contains(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestSortedSet_ContainsByComparator(t *testing.T) {
	tests := []struct {
		name string
		s    *SortedSet
		in   int
		out  bool
	}{
		{"empty", New(), 1, false},
		{"!empty/false", sortedset(10), 20, false},
		{"!empty/true", sortedset(10), 5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.ContainsByComparator(test.in, compareInt), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestSortedSet_Do(t *testing.T) {
	strResult := "P1:0 P2:0 P1:1 P2:1 P1:3 P2:3 P1:5 P2:5 "
	str := ""
	procedure1 := func(v interface{}) {
		str += fmt.Sprintf("P1:%v ", v)
	}
	procedure2 := func(v interface{}) {
		str += fmt.Sprintf("P2:%v ", v)
	}
	tests := []struct {
		name string
		s    *SortedSet
		in   []func(v interface{})
	}{
		{"empty", New(), []func(v interface{}){procedure1, procedure2}},
		{"!empty/emptyParams", NewBySlice([]interface{}{5, 3, 1, 0}, compareInt), []func(v interface{}){}},
		{"!empty", NewBySlice([]interface{}{5, 3, 1, 0}, compareInt), []func(v interface{}){procedure1, procedure2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.s.Do(test.in...)
			if test.name == "!empty" && strResult != str {
				tt.Errorf("Got: %v, Expected: %v", str, strResult)
			}
		})
	}
}
func TestSortedSet_Max(t *testing.T) {
	tests := []struct {
		name string
		s    *SortedSet
		out  interface{}
	}{
		{"empty", New(), nil},
		{"!empty", NewBySlice([]interface{}{1, 5, 6, 100, 7, 9, 4}, compareInt), 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Max(), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestSortedSet_Min(t *testing.T) {
	tests := []struct {
		name string
		s    *SortedSet
		out  interface{}
	}{
		{"empty", New(), nil},
		{"!empty", NewBySlice([]interface{}{1, 5, 6, 0, 7, 9, 4}, compareInt), 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Min(), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestSortedSet_Push(t *testing.T) {
	tests := []struct {
		name      string
		s         *SortedSet
		in        int
		toCompare []interface{}
	}{
		{"empty", New(), 0, []interface{}{0}},
		{"balanced", NewBySlice([]interface{}{12, 8, 14, 4, 10, 16, 2, 6}, compareInt), 13, []interface{}{12, 8, 14, 4, 10, 16, 2, 6, 13}},
		{"update", NewBySlice([]interface{}{12, 8, 14, 4, 10, 16, 2, 6}, compareInt), 10, []interface{}{12, 8, 14, 4, 10, 16, 2, 6}},
		{"left-left", NewBySlice([]interface{}{12, 8, 14, 4, 10, 16, 2, 6}, compareInt), 1, []interface{}{12, 8, 14, 4, 10, 16, 2, 6, 1}},
		{"right-right", NewBySlice([]interface{}{4, 2, 8, 6, 10}, compareInt), 12, []interface{}{4, 2, 8, 6, 10, 12}},
		{"left-right", NewBySlice([]interface{}{12, 8, 14, 4, 10, 16, 2, 6}, compareInt), 7, []interface{}{12, 8, 14, 4, 10, 16, 2, 6, 7}},
		{"right-left", NewBySlice([]interface{}{10, 4, 14, 2, 8, 12, 16, 6, 18}, compareInt), 17, []interface{}{10, 4, 14, 2, 8, 12, 16, 6, 18, 17}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.s.Push(test.in, compareInt)
			if !checkOrder(test.s.root, compareInt) {
				tt.Errorf("checkOrder: FAIL")
			}
			if !checkHeight(test.s.root) {
				tt.Errorf("checkHeight: FAIL")
			}
			if !checkLength(test.s.root) {
				tt.Errorf("checkLength: FAIL")
			}
			if !checkValues(test.s.root, test.toCompare) {
				tt.Errorf("checkValues: FAIL")
			}
		})
	}
}
func TestSortedSet_Remove(t *testing.T) {
	tests := []struct {
		name      string
		s         *SortedSet
		in        int
		out       bool
		toCompare []interface{}
	}{
		{"empty", New(), 5, false, []interface{}{}},
		{"!empty/false", NewBySlice([]interface{}{5, 3, 8, 1}, compareInt), 21, false, []interface{}{1, 3, 5, 8}},
		{"simple", NewBySlice([]interface{}{2, 1, 4}, compareInt), 1, true, []interface{}{2, 4}},
		{"balanced", NewBySlice([]interface{}{6, 3, 8, 2, 4}, compareInt), 3, true, []interface{}{2, 4, 6, 8}},
		{"left-left", NewBySlice([]interface{}{8, 4, 9, 2, 6}, compareInt), 9, true, []interface{}{2, 4, 6, 8}},
		{"right-right", NewBySlice([]interface{}{2, 1, 6, 4, 8}, compareInt), 1, true, []interface{}{2, 4, 6, 8}},
		{"left-right", NewBySlice([]interface{}{6, 2, 7, 4}, compareInt), 7, true, []interface{}{2, 4, 6}},
		{"right-left", NewBySlice([]interface{}{1, 2, 6, 4}, compareInt), 1, true, []interface{}{2, 4, 6}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Remove(test.in, compareInt), test.out; got != expected {
				tt.Errorf("SortedSet.Remove() -> Got: %v, Expected: %v", got, expected)
			}
			slice := test.s.Slice()
			if !checkValues(test.s.root, test.toCompare) {
				tt.Errorf("checkValues: FAIL")
			}
			if len(slice) != len(test.toCompare) {
				tt.Errorf("length: FAIL")
			}
		})
	}
}
func TestSortedSet_Slice(t *testing.T) {
	tests := []struct {
		name string
		s    *SortedSet
		out  []int
	}{
		{"empty", New(), make([]int, 0, 0)},
		{"!empty", NewBySlice([]interface{}{5, 3, 8, 4, 1}, compareInt), []int{1, 3, 4, 5, 8}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			slice := test.s.Slice()
			for i, got := range slice {
				if expected := test.out[i]; got != expected {
					tt.Errorf("Got: %v, Expected: %v", got, expected)
				}
			}
			if cap(slice) != cap(test.out) {
				tt.Errorf("capacity: FAIL")
			}
			if len(slice) != len(test.out) {
				tt.Errorf("length: FAIL")
			}
		})
	}
}
