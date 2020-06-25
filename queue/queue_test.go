// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package queue

import (
	"fmt"
	"testing"
)

func checkValuesAndOrder(q *Queue, values []interface{}) bool {
	if q.Len() != len(values) {
		return false
	}
	if q.IsEmpty() {
		return true
	}
	for n, i := q.front, 0; n != nil; n, i = n.next, i+1 {
		if n.value != values[i] {
			return false
		}
	}
	return q.back.value == values[len(values)-1]
}
func checkZeroValue(q *Queue) bool {
	return q.front == nil && q.back == nil && q.len == 0
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
		name string
		in   []interface{}
	}{
		{"empty", []interface{}{}},
		{"!empty", []interface{}{5, 2, 1, 10, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			q := NewBySlice(test.in)
			if !checkValuesAndOrder(q, test.in) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}

func TestQueue_Clone(t *testing.T) {
	tests := []struct {
		name   string
		q      *Queue
		values []interface{}
	}{
		{"empty", New(), []interface{}{}},
		{"!empty", NewBySlice([]interface{}{5, 2, 3, 6, 8, 1, 10, 4}), []interface{}{5, 2, 3, 6, 8, 1, 10, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			clone := test.q.Clone()
			if !checkValuesAndOrder(clone, test.values) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_Do(t *testing.T) {
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
		q    *Queue
		in   []func(v interface{})
	}{
		{"empty", New(), []func(v interface{}){procedure1, procedure2}},
		{"!empty/emptyParams", NewBySlice([]interface{}{0, 1, 3, 5}), []func(v interface{}){}},
		{"!empty", NewBySlice([]interface{}{0, 1, 3, 5}), []func(v interface{}){procedure1, procedure2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.q.Do(test.in...)
			if test.name == "!empty" && strResult != str {
				tt.Errorf("Got: %v, Expected: %v", str, strResult)
			}
			if !checkZeroValue(test.q) {
				t.Errorf("checkZeroValue: FAIL")
			}
		})
	}
}
func TestQueue_Equals(t *testing.T) {
	tests := []struct {
		name string
		q    *Queue
		in   *Queue
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
			if got, expected := test.q.Equals(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_EqualsByComparator(t *testing.T) {
	tests := []struct {
		name string
		q    *Queue
		in   *Queue
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

			if got, expected := test.q.EqualsByComparator(test.in, equalsInt), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_Get(t *testing.T) {
	tests := []struct {
		name      string
		q         *Queue
		out       interface{}
		toCompare []interface{}
	}{
		{"empty", New(), nil, []interface{}{}},
		{"!empty/true", NewBySlice([]interface{}{0, 1, 2}), 0, []interface{}{1, 2}},
		{"!empty/true", NewBySlice([]interface{}{0}), 0, []interface{}{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.q.Get(), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
			if !checkValuesAndOrder(test.q, test.toCompare) {
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
		q         *Queue
		out       []interface{}
		toCompare []interface{}
	}{
		{"empty", New(), nil, []interface{}{}},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2, 3}), []interface{}{}, []interface{}{0, 1, 2, 3}},
		{"!empty/true", NewBySlice([]interface{}{2, 1, 0, 2}), []interface{}{2, 1}, []interface{}{0, 2}},
		{"!empty/true", NewBySlice([]interface{}{2}), []interface{}{2}, []interface{}{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			got := test.q.GetIf(condition)
			for i, v := range got {
				if v != test.out[i] {
					tt.Errorf("Got: %v, Expected: %v", v, test.out[i])
				}
			}
			if !checkValuesAndOrder(test.q, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_Push(t *testing.T) {
	tests := []struct {
		name      string
		q         *Queue
		in        int
		toCompare []interface{}
	}{
		{"empty", New(), 5, []interface{}{5}},
		{"!empty", NewBySlice([]interface{}{0, 1, 2}), 5, []interface{}{0, 1, 2, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.q.Push(test.in)
			if !checkValuesAndOrder(test.q, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestQueue_Search(t *testing.T) {
	tests := []struct {
		name string
		q    *Queue
		in   interface{}
		out  int
	}{
		{"empty", New(), 5, -1},
		{"!empty/!match", NewBySlice([]interface{}{0, 1, 2}), 5, -1},
		{"!empty/match", NewBySlice([]interface{}{0, 1, 2}), 1, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.q.Search(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_SearchByComparator(t *testing.T) {
	tests := []struct {
		name string
		q    *Queue
		in   interface{}
		out  int
	}{
		{"empty", New(), 5, -1},
		{"!empty/!match", NewBySlice([]interface{}{0, 1, 2}), 5, -1},
		{"!empty/match", NewBySlice([]interface{}{0, 1, 2}), 1, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.q.SearchByComparator(test.in, equalsInt), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestQueue_Slice(t *testing.T) {
	tests := []struct {
		name string
		q    *Queue
		out  []interface{}
	}{
		{"empty", New(), make([]interface{}, 0, 0)},
		{"!empty", NewBySlice([]interface{}{5, 3, 8, 4, 1}), []interface{}{5, 3, 8, 4, 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			slice := test.q.Slice()
			if !checkValuesAndOrder(test.q, slice) {
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
