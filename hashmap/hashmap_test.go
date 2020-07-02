// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package hashmap

import (
	"fmt"
	coll "github.com/maguerrido/collection"
	"testing"
)

type pair struct {
	bucket int
	key    coll.Hashable
	value  interface{}
}

type key struct {
	i int
}

func (k key) Equals(v coll.Hashable) bool {
	val, ok := v.(key)
	return ok && k.i == val.i
}
func (k key) Hash() int {
	return k.i
}

func buckets(cap int, values []pair) [][]pair {
	buckets := make([][]pair, cap, cap)
	for _, v := range values {
		buckets[v.bucket] = append(buckets[v.bucket], v)
	}
	return buckets
}
func checkBucket(n *node, pairs []pair) bool {
	i := 0
	for ; n != nil && i < len(pairs); n, i = n.next, i+1 {
		if n.key != pairs[i].key || n.value != pairs[i].value {
			return false
		}
	}
	return n == nil && i == len(pairs)
}
func checkBuckets(buckets []*node, pairs [][]pair) bool {
	for i, b := range buckets {
		if !checkBucket(b, pairs[i]) {
			return false
		}
	}
	return true
}
func hmByMap(values map[coll.Hashable]interface{}) *HashMap {
	hm := New(DefaultCapacity, DefaultLoadFactor)
	for k, v := range values {
		hm.Push(k, v)
	}
	return hm
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		cap        int
		loadFactor float64
		out        *HashMap
	}{
		{"default_cap", 0, 0.30,
			&HashMap{
				buckets:    make([]*node, DefaultCapacity, DefaultCapacity),
				cap:        DefaultCapacity,
				len:        0,
				loadFactor: 0.30,
			}},
		{"default_loadFactor", 20, 0,
			&HashMap{
				buckets:    make([]*node, 20, 20),
				cap:        20,
				len:        0,
				loadFactor: DefaultLoadFactor,
			}},
		{"common", 20, 0.80,
			&HashMap{
				buckets:    make([]*node, 20, 20),
				cap:        20,
				len:        0,
				loadFactor: 0.80,
			}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			hp := New(test.cap, test.loadFactor)
			if got, expected := hp.cap, test.out.cap; got != expected {
				tt.Errorf("Capacity: Got: %v, Expected: %v", got, expected)
			}
			if got, expected := hp.len, test.out.len; got != expected {
				tt.Errorf("Length: Got: %v, Expected: %v", got, expected)
			}
			if got, expected := hp.loadFactor, test.out.loadFactor; got != expected {
				tt.Errorf("LoadFactor: Got: %v, Expected: %v", got, expected)
			}
			if got, expected := cap(hp.buckets), cap(test.out.buckets); got != expected {
				tt.Errorf("Buckets capacity: Got: %v, Expected: %v", got, expected)
			}
			if got, expected := len(hp.buckets), len(test.out.buckets); got != expected {
				tt.Errorf("Buckets capacity: Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestNewByMap(t *testing.T) {
	tests := []struct {
		name    string
		in      map[coll.Hashable]interface{}
		buckets [][]pair
	}{
		{"empty", map[coll.Hashable]interface{}{}, buckets(DefaultCapacity, []pair{})},
		{"!empty",
			map[coll.Hashable]interface{}{
				key{0}: 0,
				key{1}: 1,
				key{2}: 2,
				key{3}: 3,
			},
			buckets(DefaultCapacity, []pair{
				{0, key{0}, 0},
				{1, key{1}, 1},
				{2, key{2}, 2},
				{3, key{3}, 3},
			})},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			hm := NewByMap(test.in, DefaultCapacity, DefaultLoadFactor)
			if !checkBuckets(hm.buckets, test.buckets) {
				tt.Errorf("checkBuckets: FAIL")
			}
		})
	}
}

func TestHashMap_Clone(t *testing.T) {
	tests := []struct {
		name    string
		hm      *HashMap
		buckets [][]pair
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), buckets(DefaultCapacity, []pair{})},
		{"!empty", NewByMap(map[coll.Hashable]interface{}{
			key{0}: 0,
			key{1}: 1,
			key{2}: 2,
		}, 20, 0.70), buckets(20, []pair{
			{0, key{0}, 0},
			{1, key{1}, 1},
			{2, key{2}, 2},
		})},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			clone := test.hm.Clone()
			if !checkBuckets(clone.buckets, test.buckets) {
				tt.Errorf("checkBuckets: FAIL")
			}
			if got, expected := clone.len, test.hm.len; got != expected {
				tt.Errorf("Got: %v, Epected: %v", got, expected)
			}
			if got, expected := clone.cap, test.hm.cap; got != expected {
				tt.Errorf("Got: %v, Epected: %v", got, expected)
			}
			if got, expected := clone.loadFactor, test.hm.loadFactor; got != expected {
				tt.Errorf("Got: %v, Epected: %v", got, expected)
			}
		})
	}
}
func TestHashMap_Do(t *testing.T) {
	strResult := "P1:0 P2:0 P1:1 P2:1 P1:2 P2:2 P1:3 P2:3 "
	str := ""
	procedure1 := func(v interface{}) {
		str += fmt.Sprintf("P1:%v ", v)
	}
	procedure2 := func(v interface{}) {
		str += fmt.Sprintf("P2:%v ", v)
	}
	tests := []struct {
		name string
		hm   *HashMap
		in   []func(v interface{})
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), []func(v interface{}){procedure1, procedure2}},
		{"!empty/emptyParams", NewByMap(map[coll.Hashable]interface{}{
			key{0}: 0,
			key{1}: 1,
			key{2}: 2,
			key{3}: 3,
		}, DefaultCapacity, DefaultLoadFactor), []func(v interface{}){}},
		{"!empty", NewByMap(map[coll.Hashable]interface{}{
			key{0}: 0,
			key{1}: 1,
			key{2}: 2,
			key{3}: 3,
		}, DefaultCapacity, DefaultLoadFactor), []func(v interface{}){procedure1, procedure2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.hm.Do(test.in...)
			switch test.name {
			case "empty", "!empty/emptyParams":
				if str != "" {
					tt.Errorf("Do: FAIL")
				}
			case "!empty":
				if str != strResult {
					tt.Errorf("Do: FAIL")
				}
			}
			str = ""
		})
	}
}
func TestHashMap_Get(t *testing.T) {
	tests := []struct {
		name  string
		hm    *HashMap
		in    key
		vOut  interface{}
		okOut bool
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), key{0}, nil, false},
		{"empty/false", NewByMap(map[coll.Hashable]interface{}{
			key{0}: 0,
		}, DefaultCapacity, DefaultLoadFactor), key{1}, nil, false},
		{"empty/true", NewByMap(map[coll.Hashable]interface{}{
			key{0}:  0,
			key{15}: 15,
		}, DefaultCapacity, DefaultLoadFactor), key{15}, 15, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			vGot, okGot := test.hm.Get(test.in)
			vExpected := test.vOut
			okExpected := test.okOut
			if vGot != vExpected {
				tt.Errorf("Got: %v, Expected: %v", vGot, vExpected)
			}
			if okGot != okExpected {
				tt.Errorf("Got: %v, Expected: %v", okGot, okExpected)
			}
		})
	}
}
func TestHashMap_Map(t *testing.T) {
	tests := []struct {
		name string
		hm   *HashMap
		out  map[coll.Hashable]interface{}
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), map[coll.Hashable]interface{}{}},
		{"!empty", NewByMap(map[coll.Hashable]interface{}{
			key{0}:  0,
			key{5}:  5,
			key{16}: 16,
		}, DefaultCapacity, DefaultLoadFactor), map[coll.Hashable]interface{}{
			key{0}:  0,
			key{16}: 16,
			key{5}:  5,
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := test.hm.Map()
			if got, expected := len(m), len(test.out); got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
			for k, got := range m {
				if expected := test.out[k]; got != expected {
					tt.Errorf("Got: %v, Expected: %v", got, expected)
				}
			}
		})
	}
}
func TestHashMap_Push(t *testing.T) {
	tests := []struct {
		name  string
		key   coll.Hashable
		value int
		out   bool
		hm    *HashMap
		pairs [][]pair
	}{
		{"empty", key{0}, 0, true,
			New(DefaultCapacity, DefaultLoadFactor),
			buckets(DefaultCapacity, []pair{
				{0, key{0}, 0},
			})},
		{"empty/nilKey", nil, 0, false,
			New(DefaultCapacity, DefaultLoadFactor),
			buckets(DefaultCapacity, []pair{})},
		{"!empty/different_bucket", key{0}, 0, true,
			hmByMap(map[coll.Hashable]interface{}{
				key{1}: 1}),
			buckets(DefaultCapacity, []pair{
				{1, key{1}, 1},
				{0, key{0}, 0},
			})},
		{"!empty/same_bucket", key{0}, 0, true,
			hmByMap(map[coll.Hashable]interface{}{
				key{16}: 16}),
			buckets(DefaultCapacity, []pair{
				{0, key{0}, 0},
				{0, key{16}, 16},
			})},
		{"!empty/update", key{0}, 0, true,
			hmByMap(map[coll.Hashable]interface{}{
				key{0}: 16}), // 16 should be updated by 0
			buckets(DefaultCapacity, []pair{
				{0, key{0}, 0},
			})},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.hm.Push(test.key, test.value), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
			if !checkBuckets(test.hm.buckets, test.pairs) {
				tt.Errorf("checkBuckets: FAIL")
			}
		})
	}
}
func TestHashMap_PushReHashing(t *testing.T) {
	hm := NewByMap(map[coll.Hashable]interface{}{
		key{0}:  0,
		key{1}:  1,
		key{2}:  2,
		key{3}:  3,
		key{4}:  4,
		key{5}:  5,
		key{6}:  6,
		key{7}:  7,
		key{8}:  8,
		key{9}:  9,
		key{10}: 10,
		key{11}: 11,
		key{12}: 12,
	}, DefaultCapacity, DefaultLoadFactor)
	inserts := 13
	capGot, capExpected := hm.cap, DefaultCapacity
	loadFactorGot, loadFactorExpected := hm.loadFactor, DefaultLoadFactor
	lenGot, lenExpected := hm.len, inserts
	lenBucketsGot, lenBucketsExpected := len(hm.buckets), DefaultCapacity
	if capGot != capExpected ||
		loadFactorGot != loadFactorExpected ||
		lenGot != lenExpected ||
		lenBucketsGot != lenBucketsExpected {
		t.Errorf("Fail before reHashing")
		t.Errorf("CapGot: %v, CapExpected: %v\n", capGot, capExpected)
		t.Errorf("LoadFactorGot: %v, LoadFactorExpected: %v\n", loadFactorGot, loadFactorExpected)
		t.Errorf("LenGot: %v, LenExpected: %v\n", lenGot, lenExpected)
		t.Errorf("LenBucketsGot: %v, LenBucketsExpected: %v\n", lenBucketsGot, lenBucketsExpected)
	}
	hm.Push(key{13}, 13)
	capGot, capExpected = hm.cap, DefaultCapacity*2
	loadFactorGot, loadFactorExpected = hm.loadFactor, DefaultLoadFactor
	lenGot, lenExpected = hm.len, inserts+1
	lenBucketsGot, lenBucketsExpected = len(hm.buckets), DefaultCapacity*2
	if capGot != capExpected ||
		loadFactorGot != loadFactorExpected ||
		lenGot != lenExpected ||
		lenBucketsGot != lenBucketsExpected {
		t.Errorf("Fail after reHashing")
		t.Errorf("CapGot: %v, CapExpected: %v\n", capGot, capExpected)
		t.Errorf("LoadFactorGot: %v, LoadFactorExpected: %v\n", loadFactorGot, loadFactorExpected)
		t.Errorf("LenGot: %v, LenExpected: %v\n", lenGot, lenExpected)
		t.Errorf("LenBucketsGot: %v, LenBucketsExpected: %v\n", lenBucketsGot, lenBucketsExpected)
	}

}
func TestHashMap_Remove(t *testing.T) {
	tests := []struct {
		name    string
		hm      *HashMap
		in      coll.Hashable
		vOut    interface{}
		okOut   bool
		buckets [][]pair
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), key{0}, nil, false, buckets(DefaultCapacity, []pair{})},
		{"!empty/nilKey", NewByMap(map[coll.Hashable]interface{}{
			key{1}: 1,
		}, DefaultCapacity, DefaultLoadFactor), nil, nil, false, buckets(DefaultCapacity, []pair{
			{1, key{1}, 1},
		})},
		{"!empty/false", NewByMap(map[coll.Hashable]interface{}{
			key{1}:  1,
			key{16}: 16,
		}, DefaultCapacity, DefaultLoadFactor), key{0}, nil, false, buckets(DefaultCapacity, []pair{
			{0, key{16}, 16},
			{1, key{1}, 1},
		})},
		{"!empty/true/middle", nil, key{0}, 0, true, buckets(DefaultCapacity, []pair{
			{0, key{16}, 16},
			{0, key{32}, 32},
			{14, key{14}, 14},
		})},
		{"!empty/true/front", nil, key{0}, 0, true, buckets(DefaultCapacity, []pair{
			{0, key{16}, 16},
			{0, key{32}, 32},
			{14, key{14}, 14},
		})},
		{"!empty/true/back", nil, key{0}, 0, true, buckets(DefaultCapacity, []pair{
			{0, key{16}, 16},
			{0, key{32}, 32},
			{14, key{14}, 14},
		})},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			switch test.name {
			case "!empty/true/middle":
				test.hm = New(DefaultCapacity, DefaultLoadFactor)
				test.hm.Push(key{32}, 32)
				test.hm.Push(key{0}, 0)
				test.hm.Push(key{16}, 16)
				test.hm.Push(key{14}, 14)
			case "!empty/true/front":
				test.hm = New(DefaultCapacity, DefaultLoadFactor)
				test.hm.Push(key{32}, 32)
				test.hm.Push(key{16}, 16)
				test.hm.Push(key{0}, 0)
				test.hm.Push(key{14}, 14)
			case "!empty/true/back":
				test.hm = New(DefaultCapacity, DefaultLoadFactor)
				test.hm.Push(key{0}, 0)
				test.hm.Push(key{32}, 32)
				test.hm.Push(key{16}, 16)
				test.hm.Push(key{14}, 14)
			}
			vGot, okGot := test.hm.Remove(test.in)
			vExpected := test.vOut
			okExpected := test.okOut
			if vGot != vExpected {
				tt.Errorf("Got: %v, Expected: %v", vGot, vExpected)
			}
			if okGot != okExpected {
				tt.Errorf("Got: %v, Expected: %v", okGot, okExpected)
			}
			if !checkBuckets(test.hm.buckets, test.buckets) {
				tt.Errorf("checkBuckets: FAIL")
			}
		})
	}
}
func TestHashMap_Search(t *testing.T) {
	tests := []struct {
		name string
		hm   *HashMap
		in   interface{}
		out  coll.Hashable
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), 0, nil},
		{"!empty/false", NewByMap(map[coll.Hashable]interface{}{
			key{1}: 1,
		}, DefaultCapacity, DefaultLoadFactor), 0, nil},
		{"!empty/true", NewByMap(map[coll.Hashable]interface{}{
			key{13}: 13,
		}, DefaultCapacity, DefaultLoadFactor), 13, key{13}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.hm.Search(test.in), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
func TestHashMap_SearchByComparator(t *testing.T) {
	tests := []struct {
		name string
		hm   *HashMap
		in   interface{}
		out  coll.Hashable
	}{
		{"empty", New(DefaultCapacity, DefaultLoadFactor), 0, nil},
		{"!empty/false", NewByMap(map[coll.Hashable]interface{}{
			key{1}: 1,
		}, DefaultCapacity, DefaultLoadFactor), 0, nil},
		{"!empty/true", NewByMap(map[coll.Hashable]interface{}{
			key{13}: 13,
		}, DefaultCapacity, DefaultLoadFactor), 13, key{13}},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if got, expected := test.hm.SearchByComparator(test.in, func(v1, v2 interface{}) bool {
				i1, ok1 := v1.(int)
				i2, ok2 := v2.(int)
				return ok1 && ok2 && i1 == i2
			}), test.out; got != expected {
				tt.Errorf("Got: %v, Expected: %v", got, expected)
			}
		})
	}
}
