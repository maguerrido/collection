// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package list

import (
	"fmt"
	coll "github.com/maguerrido/collection"
	"testing"
)

func checkValuesAndOrder(l *List, values []interface{}) bool {
	if l.Len() != len(values) {
		return false
	}
	if l.IsEmpty() {
		return true
	}
	for e, i := l.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		if e.Value() != values[i] {
			return false
		}
	}
	for e, i := l.Back(), len(values)-1; e != nil; e, i = e.Prev(), i-1 {
		if e.Value() != values[i] {
			return false
		}
	}
	return true

}
func checkZeroValue(l *List) bool {
	return l.Front() == nil && l.Back() == nil && l.Len() == 0
}
func compareInt(v1, v2 interface{}) int {
	int1 := v1.(int)
	int2 := v2.(int)
	return int1 - int2
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
			l := NewBySlice(test.in)
			if !checkValuesAndOrder(l, test.in) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}

func TestList_Clone(t *testing.T) {
	tests := []struct {
		name   string
		l      *List
		values []interface{}
	}{
		{"empty", New(), []interface{}{}},
		{"!empty", NewBySlice([]interface{}{5, 2, 3, 6, 8, 1, 10, 4}), []interface{}{5, 2, 3, 6, 8, 1, 10, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			clone := test.l.Clone()
			if !checkValuesAndOrder(clone, test.values) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_Do(t *testing.T) {
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
		l    *List
		in   []func(v interface{})
	}{
		{"empty", New(), []func(v interface{}){procedure1, procedure2}},
		{"!empty/emptyParams", NewBySlice([]interface{}{0, 1, 3, 5}), []func(v interface{}){}},
		{"!empty", NewBySlice([]interface{}{0, 1, 3, 5}), []func(v interface{}){procedure1, procedure2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.Do(test.in...)
			if test.name == "!empty" && strResult != str {
				tt.Errorf("Got: %v, Expected: %v", str, strResult)
			}
		})
	}
}
func TestList_Equals(t *testing.T) {
	tests := []struct {
		name string
		l    *List
		in   *List
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
			if got, expected := test.l.Equals(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestList_EqualsByComparator(t *testing.T) {
	tests := []struct {
		name string
		l    *List
		in   *List
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

			if got, expected := test.l.EqualsByComparator(test.in, equalsInt), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestList_Get(t *testing.T) {
	tests := []struct {
		name     string
		l        *List
		in       int
		out      interface{}
		contains bool
	}{
		{"empty", New(), 0, nil, false},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2}), 3, nil, false},
		{"!empty/true", NewBySlice([]interface{}{0, 1, 2}), 1, 1, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			got := test.l.Get(test.in)
			if test.l.Contains(got) != test.contains {
				tt.Errorf("Get: FAIL")
			}
			if test.contains {
				if got, expected := got.Value(), test.out; got != expected {
					tt.Errorf("Got: %v, Expected: %v", got, expected)
				}
			}

		})
	}
}
func TestList_MoveAfter(t *testing.T) {
	t.Run("empty/false", func(tt *testing.T) {
		l := New()
		l2 := NewBySlice([]interface{}{0, 1, 2})
		e := l2.Get(0)
		mark := l2.Get(1)
		if l.MoveAfter(e, mark) {
			tt.Errorf("MoveAfter: FAIL")
		}
	})
	t.Run("!empty/false", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		l2 := NewBySlice([]interface{}{0, 1, 2})
		e := l2.Get(0)
		mark := l2.Get(1)
		if l.MoveAfter(e, mark) {
			tt.Errorf("MoveAfter: FAIL")
		}
	})
	t.Run("!empty/true/two-elements", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1})
		e := l.Get(0)
		mark := l.Get(1)
		if !l.MoveAfter(e, mark) {
			tt.Errorf("List.MoveAfter() -> False: Fail")
		}
		if !checkValuesAndOrder(l, []interface{}{1, 0}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/e==front,mark==back", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		e := l.Get(0)
		mark := l.Get(2)
		if !l.MoveAfter(e, mark) {
			tt.Errorf("MoveAfter: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{1, 2, 0}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/e==back", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		e := l.Get(2)
		mark := l.Get(0)
		if !l.MoveAfter(e, mark) {
			tt.Errorf("MoveAfter: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 2, 1}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/no-changes", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2, 3})
		e := l.Get(2)
		mark := l.Get(1)
		if !l.MoveAfter(e, mark) {
			tt.Errorf("MoveAfter: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 1, 2, 3}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/middle", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2, 3})
		e := l.Get(1)
		mark := l.Get(2)
		if !l.MoveAfter(e, mark) {
			tt.Errorf("MoveAfter: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 2, 1, 3}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
}
func TestList_MoveBefore(t *testing.T) {
	t.Run("empty/false", func(tt *testing.T) {
		l := New()
		l2 := NewBySlice([]interface{}{0, 1, 2})
		e := l2.Get(0)
		mark := l2.Get(1)
		if l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
	})
	t.Run("!empty/false", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		l2 := NewBySlice([]interface{}{0, 1, 2})
		e := l2.Get(0)
		mark := l2.Get(1)
		if l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
	})
	t.Run("!empty/true/two-elements", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1})
		e := l.Get(1)
		mark := l.Get(0)
		if !l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{1, 0}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/e==back,mark==front", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		e := l.Get(2)
		mark := l.Get(0)
		if !l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{2, 0, 1}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/e==front", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		e := l.Get(0)
		mark := l.Get(2)
		if !l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{1, 0, 2}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/no-changes", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2, 3})
		e := l.Get(1)
		mark := l.Get(2)
		if !l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 1, 2, 3}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/middle", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2, 3})
		e := l.Get(2)
		mark := l.Get(1)
		if !l.MoveBefore(e, mark) {
			tt.Errorf("MoveBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 2, 1, 3}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
}
func TestList_PushAfter(t *testing.T) {
	t.Run("empty/false", func(tt *testing.T) {
		l := New()
		l2 := NewBySlice([]interface{}{0})
		mark := l2.Get(0)
		if got := l.PushAfter(5, mark); got != nil {
			tt.Errorf("PushAfter: FAIL")
		}
	})
	t.Run("!empty/false", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		l2 := NewBySlice([]interface{}{0})
		mark := l2.Get(0)
		if got := l.PushAfter(5, mark); got != nil {
			tt.Errorf("PushAfter: FAIL")
		}
	})
	t.Run("!empty/true/back", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0})
		mark := l.Get(0)
		v := 5
		if got := l.PushAfter(v, mark); got == nil || got.Value() != v {
			tt.Errorf("PushAfter: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 5}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/middle", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1})
		mark := l.Get(0)
		v := 5
		if got := l.PushAfter(v, mark); got == nil || got.Value() != v {
			tt.Errorf("PushAfter: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 5, 1}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
}
func TestList_PushBack(t *testing.T) {
	tests := []struct {
		name      string
		l         *List
		in        int
		toCompare []interface{}
	}{
		{"empty", New(), 5, []interface{}{5}},
		{"!empty", NewBySlice([]interface{}{0, 1, 2}), 5, []interface{}{0, 1, 2, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.PushBack(test.in)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_PushBackList(t *testing.T) {
	tests := []struct {
		name      string
		l         *List
		in        *List
		toCompare []interface{}
	}{
		{"empty/empty", New(), New(), []interface{}{}},
		{"empty/!empty", New(), NewBySlice([]interface{}{0, 1, 2}), []interface{}{0, 1, 2}},
		{"!empty/empty", NewBySlice([]interface{}{0, 1, 2}), New(), []interface{}{0, 1, 2}},
		{"!empty/!empty", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{3, 4, 5}), []interface{}{0, 1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.PushBackList(test.in)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_PushBefore(t *testing.T) {
	t.Run("empty/false", func(tt *testing.T) {
		l := New()
		l2 := NewBySlice([]interface{}{0})
		mark := l2.Get(0)
		if got := l.PushBefore(5, mark); got != nil {
			tt.Errorf("PushBefore: FAIL")
		}
	})
	t.Run("!empty/false", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		l2 := NewBySlice([]interface{}{0})
		mark := l2.Get(0)
		if got := l.PushBefore(5, mark); got != nil {
			tt.Errorf("PushBefore: FAIL")
		}
	})
	t.Run("!empty/true/front", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0})
		mark := l.Get(0)
		v := 5
		if got := l.PushBefore(v, mark); got == nil || got.Value() != v {
			tt.Errorf("PushBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{5, 0}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/true/middle", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1})
		mark := l.Get(1)
		v := 5
		if got := l.PushBefore(v, mark); got == nil || got.Value() != v {
			tt.Errorf("PushBefore: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 5, 1}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
}
func TestList_PushFront(t *testing.T) {
	tests := []struct {
		name      string
		l         *List
		in        int
		toCompare []interface{}
	}{
		{"empty", New(), 5, []interface{}{5}},
		{"!empty", NewBySlice([]interface{}{0, 1, 2}), 5, []interface{}{5, 0, 1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.PushFront(test.in)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_PushFrontList(t *testing.T) {
	tests := []struct {
		name      string
		l         *List
		in        *List
		toCompare []interface{}
	}{
		{"empty/empty", New(), New(), []interface{}{}},
		{"empty/!empty", New(), NewBySlice([]interface{}{0, 1, 2}), []interface{}{0, 1, 2}},
		{"!empty/empty", NewBySlice([]interface{}{0, 1, 2}), New(), []interface{}{0, 1, 2}},
		{"!empty/!empty", NewBySlice([]interface{}{0, 1, 2}), NewBySlice([]interface{}{3, 4, 5}), []interface{}{3, 4, 5, 0, 1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.PushFrontList(test.in)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_Remove(t *testing.T) {
	tests := []struct {
		name      string
		l         *List
		in        int
		toCompare []interface{}
	}{
		{"empty", New(), 5, []interface{}{}},
		{"!empty/false", NewBySlice([]interface{}{0, 1, 2}), 3, []interface{}{0, 1, 2}},
		{"!empty/true", NewBySlice([]interface{}{0, 1, 1, 2}), 1, []interface{}{0, 1, 2}},
		{"!empty/true", NewBySlice([]interface{}{0}), 0, []interface{}{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.Remove(test.in)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_RemoveElement(t *testing.T) {
	t.Run("empty", func(tt *testing.T) {
		l := New()
		v, ok := l.RemoveElement(l.Front())
		if v != nil || ok {
			tt.Errorf("RemoveElement: FAIL")
		}
	})
	t.Run("!empty/false", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		l2 := NewBySlice([]interface{}{0, 1, 2})
		v, ok := l.RemoveElement(l2.Front())
		if v != nil || ok {
			tt.Errorf("RemoveElement: FAIL")
		}
	})
	t.Run("!empty/true", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		v, ok := l.RemoveElement(l.Front())
		if v == nil || v != 0 || !ok {
			tt.Errorf("RemoveElement: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{1, 2}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
}
func TestList_RemoveFrom(t *testing.T) {
	t.Run("empty", func(tt *testing.T) {
		l := New()
		n := l.RemoveFrom(l.Front(), l.Back())
		if n != 0 {
			tt.Errorf("RemoveFrom: FAIL")
		}
	})
	t.Run("!empty/0", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2})
		l2 := NewBySlice([]interface{}{0, 1, 2})
		n := l.RemoveFrom(l2.Front(), l2.Back())
		if n != 0 {
			tt.Errorf("RemoveFrom: FAIL")
		}
	})
	t.Run("!empty/start-to-back", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2, 3})
		_, start := l.Search(1)
		_, end := l.Search(0)
		n := l.RemoveFrom(start, end)
		if n != 3 {
			tt.Errorf("RemoveFrom: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
	t.Run("!empty/start-to-end", func(tt *testing.T) {
		l := NewBySlice([]interface{}{0, 1, 2, 3})
		_, start := l.Search(1)
		_, end := l.Search(2)
		n := l.RemoveFrom(start, end)
		if n != 2 {
			tt.Errorf("RemoveFrom: FAIL")
		}
		if !checkValuesAndOrder(l, []interface{}{0, 3}) {
			tt.Errorf("checkValuesAndOrder: FAIL")
		}
	})
}
func TestList_RemoveIf(t *testing.T) {
	condition := func(v interface{}) bool {
		intV := v.(int)
		return intV > 2
	}
	tests := []struct {
		name      string
		l         *List
		out       int
		toCompare []interface{}
	}{
		{"empty", New(), 0, []interface{}{}},
		{"!empty/0", NewBySlice([]interface{}{0, 1, 2}), 0, []interface{}{0, 1, 2}},
		{"!empty/>0", NewBySlice([]interface{}{0, 1, 2, 3, 4, 5}), 3, []interface{}{0, 1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			if got, expected := test.l.RemoveIf(condition), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestList_Search(t *testing.T) {
	tests := []struct {
		name     string
		l        *List
		in       int
		outIndex int
	}{
		{"empty", New(), 5, -1},
		{"!empty/!match", NewBySlice([]interface{}{0, 1, 2}), 5, -1},
		{"!empty/match", NewBySlice([]interface{}{0, 1, 2}), 1, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			i, e := test.l.Search(test.in)
			if test.outIndex == -1 {
				if i != test.outIndex || e != nil {
					tt.Errorf("Search: FAIL")
				}
			} else {
				if i != test.outIndex || e == nil || e.Value() != test.in {
					tt.Errorf("Search: FAIL")
				}
			}
		})
	}
}
func TestList_SearchByComparator(t *testing.T) {
	tests := []struct {
		name     string
		l        *List
		in       int
		outIndex int
	}{
		{"empty", New(), 5, -1},
		{"!empty/!match", NewBySlice([]interface{}{0, 1, 2}), 5, -1},
		{"!empty/match", NewBySlice([]interface{}{0, 1, 2}), 1, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			i, e := test.l.SearchByComparator(test.in, equalsInt)
			if test.outIndex == -1 {
				if i != test.outIndex || e != nil {
					tt.Errorf("Search: FAIL")
				}
			} else {
				if i != test.outIndex || e == nil || e.Value() != test.in {
					tt.Errorf("Search: FAIL")
				}
			}
		})
	}
}
func TestList_Slice(t *testing.T) {
	tests := []struct {
		name string
		l    *List
		out  []interface{}
	}{
		{"empty", New(), make([]interface{}, 0, 0)},
		{"!empty", NewBySlice([]interface{}{5, 3, 8, 4, 1}), []interface{}{5, 3, 8, 4, 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			slice := test.l.Slice()
			if !checkValuesAndOrder(test.l, slice) {
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
func TestList_Sort(t *testing.T) {
	tests := []struct {
		name      string
		l         *List
		toCompare []interface{}
	}{
		{"empty", New(), make([]interface{}, 0, 0)},
		{"selectionSort", NewBySlice([]interface{}{5, 3, 8, 4, 3, 1}), []interface{}{1, 3, 3, 4, 5, 8}},
		{"quickSort", NewBySlice([]interface{}{5, 3, 8, 4, 1, 8, 6, 1, 0, 10, 9}), []interface{}{0, 1, 1, 3, 4, 5, 6, 8, 8, 9, 10}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.l.Sort(compareInt)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}

func TestIterator_ForEach(t *testing.T) {
	action := func(v *interface{}) {
		intV, _ := (*v).(int)
		*v = intV * 2
	}
	tests := []struct {
		name      string
		in        func(v *interface{})
		l         *List
		toCompare []interface{}
	}{
		{"empty", action,
			New(),
			[]interface{}{}},
		{"!empty/emptyParams", nil,
			NewBySlice([]interface{}{0, 1, 2}),
			[]interface{}{0, 1, 2}},
		{"!empty", action,
			NewBySlice([]interface{}{0, 1, 2}),
			[]interface{}{0, 2, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			iterator := test.l.Iterator()
			iterator.ForEach(test.in)
			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}
		})
	}
}
func TestIterator_HasNext(t *testing.T) {
	tests := []struct {
		name     string
		loopNext int
		out      bool
		l        *List
	}{
		{"empty/false", 0, false,
			New()},
		{"!empty/false", 1, false,
			NewBySlice([]interface{}{0})},
		{"!empty/true/first", 0, true,
			NewBySlice([]interface{}{0})},
		{"!empty/true", 1, true,
			NewBySlice([]interface{}{0, 1})},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			iterator := test.l.Iterator()
			for i := 0; i < test.loopNext; i++ {
				iterator.HasNext()
				_, _ = iterator.Next()
			}
			if got, expected := iterator.HasNext(), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestIterator_Next(t *testing.T) {
	tests := []struct {
		name     string
		loopNext int
		out      interface{}
		errStr   string
		l        *List
	}{
		{"empty/error", 0, nil, coll.ErrorIteratorHasNext,
			New()},
		{"!empty/error", 1, nil, coll.ErrorIteratorHasNext,
			NewBySlice([]interface{}{0})},
		{"!empty/error/withOutHasNext", 0, nil, coll.ErrorIteratorNext,
			NewBySlice([]interface{}{0})},
		{"!empty/ok/first", 0, 0, "",
			NewBySlice([]interface{}{0})},
		{"!empty/ok", 1, 1, "",
			NewBySlice([]interface{}{0, 1})},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			iterator := test.l.Iterator()
			for i := 0; i < test.loopNext; i++ {
				iterator.HasNext()
				_, _ = iterator.Next()
			}

			if test.name != "!empty/error/withOutHasNext" {
				iterator.HasNext()
			}
			got, err := iterator.Next()

			if expected := test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}

			if test.errStr == "" {
				if err != nil {
					tt.Errorf("error detected: %v", err.Error())
				}
			} else {
				if err == nil {
					tt.Errorf("error not detected")
				} else {
					if err.Error() != test.errStr {
						tt.Errorf("wrong error")
					}
				}
			}
		})
	}
}
func TestIterator_Remove(t *testing.T) {
	tests := []struct {
		name      string
		loopNext  int
		errStr    string
		l         *List
		toCompare []interface{}
	}{
		{"empty/error", 0, coll.ErrorIteratorHasNext,
			New(),
			[]interface{}{}},
		{"!empty/error", 1, coll.ErrorIteratorHasNext,
			NewBySlice([]interface{}{0}),
			[]interface{}{0}},
		{"!empty/error/withOutNext", 0, coll.ErrorIteratorRemove,
			NewBySlice([]interface{}{0}),
			[]interface{}{0}},
		{"!empty/error/doubleRemove", 0, coll.ErrorIteratorRemove,
			NewBySlice([]interface{}{0, 1}),
			[]interface{}{1}},
		{"!empty/ok/first", 0, "",
			NewBySlice([]interface{}{0}),
			[]interface{}{}},
		{"!empty/ok", 1, "",
			NewBySlice([]interface{}{0, 1}),
			[]interface{}{0}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			iterator := test.l.Iterator()
			for i := 0; i < test.loopNext; i++ {
				iterator.HasNext()
				_, _ = iterator.Next()
			}

			iterator.HasNext()
			if test.name != "!empty/error/withOutNext" {
				_, _ = iterator.Next()
			}
			if test.name == "!empty/error/doubleRemove" {
				_ = iterator.Remove()
			}
			err := iterator.Remove()

			if !checkValuesAndOrder(test.l, test.toCompare) {
				tt.Errorf("checkValuesAndOrder: FAIL")
			}

			if test.errStr == "" {
				if err != nil {
					tt.Errorf("error detected: %v", err.Error())
				}
			} else {
				if err == nil {
					tt.Errorf("error not detected")
				} else {
					if err.Error() != test.errStr {
						tt.Errorf("wrong error")
					}
				}
			}
		})
	}
}
