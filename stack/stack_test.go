// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package stack

import (
	"fmt"
	"testing"
)

func checkValuesAndOrder(s *Stack, values []interface{}) bool {
	if s.Len() != len(values) {
		return false
	}
	if s.IsEmpty() {
		return true
	}
	for n, i := s.top, 0; n != nil; n, i = n.next, i+1 {
		if n.value != values[i] {
			return false
		}
	}
	return true
}
func checkZeroValue(s *Stack) bool {
	return s.top == nil && s.len == 0
}
func equalsInt(v1, v2 interface{}) bool {
	int1 := v1.(int)
	int2 := v2.(int)
	return int1 == int2
}

func TestNew(t *testing.T) {
	got := New()
	if !checkZeroValue(got) {
		t.Errorf("checkZeroValue: FAIL")
	}
}
func TestNewBySlice(t *testing.T) {
	tests := []struct {
		name      string
		in        []interface{}
		toCompare []interface{}
	}{
		{"empty", []interface{}{}, []interface{}{}},
		{"!empty", []interface{}{5, 2, 1, 10, 4}, []interface{}{4, 10, 1, 2, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			s := NewBySlice(test.in)
			if !checkValuesAndOrder(s, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}

func TestQueue_Clone(t *testing.T) {
	tests := []struct {
		name   string
		s      *Stack
		values []interface{}
	}{
		{"empty", New(), []interface{}{}},
		{"!empty", NewBySlice([]interface{}{5, 2, 3, 6, 8, 1, 10, 4}), []interface{}{4, 10, 1, 8, 6, 3, 2, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			clone := test.s.Clone()
			if !checkValuesAndOrder(clone, test.values) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_Do(t *testing.T) {
	strResult := "P1:5 P2:5 P1:3 P2:3 P1:1 P2:1 P1:0 P2:0 "
	str := ""
	procedure1 := func(v interface{}) {
		str += fmt.Sprintf("P1:%v ", v)
	}
	procedure2 := func(v interface{}) {
		str += fmt.Sprintf("P2:%v ", v)
	}
	tests := []struct {
		name string
		s    *Stack
		in   []func(v interface{})
	}{
		{"empty", New(), []func(v interface{}){procedure1, procedure2}},
		{"!empty/emptyParams", NewBySlice([]interface{}{0, 1, 3, 5}), []func(v interface{}){}},
		{"!empty", NewBySlice([]interface{}{0, 1, 3, 5}), []func(v interface{}){procedure1, procedure2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.s.Do(test.in...)
			if test.name == "!empty" && strResult != str {
				tt.Errorf("Got: %v, Expected: %v", str, strResult)
			}
			if !checkZeroValue(test.s) {
				t.Errorf("checkZeroValue: FAIL")
			}
		})
	}
}
func TestQueue_Equals(t *testing.T) {
	tests := []struct {
		name string
		s    *Stack
		in   *Stack
		out  bool
	}{
		{"empty/true", New(), New(), true},
		{"empty/false", New(), NewBySlice([]interface{}{0, 1, 2}), false},
		{"!empty/true", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{0, 1, 2}), true},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{0, 1, 2, 3}), false},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{0, 1, 1}), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Equals(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_EqualsByComparator(t *testing.T) {
	tests := []struct {
		name string
		s    *Stack
		in   *Stack
		out  bool
	}{
		{"empty/true", New(), New(), true},
		{"empty/false", New(), NewBySlice([]interface{}{0, 1, 2}), false},
		{"!empty/true", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{0, 1, 2}), true},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{0, 1, 2, 3}), false},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{0, 1, 1}), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			if got, expected := test.s.EqualsByComparator(test.in, equalsInt), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_Get(t *testing.T) {
	tests := []struct {
		name      string
		s         *Stack
		out       interface{}
		toCompare []interface{}
	}{
		{"empty", New(), nil, []interface{}{}},
		{"!empty/true", NewBySlice([]interface{}{0, 1, 2}), 2, []interface{}{1, 0}},
		{"!empty/true", NewBySlice([]interface{}{0}), 0, []interface{}{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Get(), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
			if !checkValuesAndOrder(test.s, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_GetIf(t *testing.T) {
	condition := func(v interface{}) bool {
		intV := v.(int)
		return intV > 0
	}
	tests := []struct {
		name      string
		s         *Stack
		out       []interface{}
		toCompare []interface{}
	}{
		{"empty", New(), nil, []interface{}{}},
		{"!empty/false", NewBySlice([]interface{}{3, 2, 1, 0}), []interface{}{}, []interface{}{0, 1, 2, 3}},
		{"!empty/true", NewBySlice([]interface{}{2, 0, 1, 2}), []interface{}{2, 1}, []interface{}{0, 2}},
		{"!empty/true", NewBySlice([]interface{}{2}), []interface{}{2}, []interface{}{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			got := test.s.GetIf(condition)
			for i, v := range got {
				if v != test.out[i] {
					tt.Errorf("Got: %v, Expected: %v", v, test.out[i])
				}
			}
			if !checkValuesAndOrder(test.s, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_Push(t *testing.T) {
	tests := []struct {
		name      string
		s         *Stack
		in        int
		toCompare []interface{}
	}{
		{"empty", New(), 5, []interface{}{5}},
		{"!empty", NewBySlice([]interface{}{0, 1, 2}), 5, []interface{}{5, 2, 1, 0}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.s.Push(test.in)
			if !checkValuesAndOrder(test.s, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_Search(t *testing.T) {
	tests := []struct {
		name string
		s    *Stack
		in   interface{}
		out  int
	}{
		{"empty", New(), 5, -1},
		{"!empty/!match", NewBySlice([]interface{}{0, 1, 2}), 5, -1},
		{"!empty/match", NewBySlice([]interface{}{0, 1, 2}), 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.Search(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_SearchByComparator(t *testing.T) {
	tests := []struct {
		name string
		s    *Stack
		in   interface{}
		out  int
	}{
		{"empty", New(), 5, -1},
		{"!empty/!match", NewBySlice([]interface{}{0, 1, 2}), 5, -1},
		{"!empty/match", NewBySlice([]interface{}{0, 1, 2}), 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.s.SearchByComparator(test.in, equalsInt), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_Slice(t *testing.T) {
	tests := []struct {
		name string
		s    *Stack
		out  []interface{}
	}{
		{"empty", New(), make([]interface{}, 0, 0)},
		{"!empty", NewBySlice([]interface{}{5, 3, 8, 4, 1}), []interface{}{1, 4, 8, 3, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			slice := test.s.Slice()
			if !checkValuesAndOrder(test.s, slice) {
				tt.Errorf("checkValuesAndOrder: FAIL")
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
