// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

// Package hashmap implements a hash table with has map behaviors.
// This implementation does not define any kind of order when saving or obtaining values.
// Collisions are handled through the separate chaining method: each bucket contains chained nodes (singly-linked list)
//related by unique hash code (generated by the hash function), that is, if two or more values return the same hash
//code, both will be stored in the same bucket.
// When the number of entries exceeds the product of the load factor and the current capacity, the hash map will double
//its capacity and all entries will be reinserted. This process is called rehashing.
// The value of the load factor and the capacity are defined when using the New constructor, once defined it cannot be
//changed manually.
// Each key of a key-value pair must implement the Hashable interface of the Collection package. This ensures that the
//keys can be compared and can generate their hash code.
package hashmap

import (
	"fmt"
	coll "github.com/maguerrido/collection"
)

const (
	// DefaultCapacity will be the capacity when the constructor receives an integer less than or equal to zero.
	DefaultCapacity = 16

	// DefaultLoadFactor will be the load factor when the constructor receives an integer less than or equal to zero.
	DefaultLoadFactor = 0.75
)

// node of a list belonging to a bucket.
type node struct {
	// hasCode is the hash generated by the key Hash method.
	hashCode int

	// key is the key of a key-value pair.
	key coll.Hashable

	// value is the value of a key-value pair.
	value interface{}

	// next points to the next node.
	// If this node is the back of the list then next points to nil.
	next *node
}

// clear sets the properties of the node to its zero values.
// Time complexity: O(1).
func (n *node) clear() {
	n.hashCode, n.key, n.value, n.next = 0, nil, nil, nil
}

// search returns the linked node it stores to 'key'.
// If the 'key' is not found, then returns nil.
// The comparison between keys is defined by its Equals method (Hashable interface).
// Time complexity: O(n), where n is the current length of the list.
func (n *node) search(key coll.Hashable) *node {
	for i := n; i != nil; i = i.next {
		if i.key.Equals(key) {
			return i
		}
	}
	return nil
}

// HashMap represents a hash table that stores key-value pairs.
// The zero value of HashMap is NOT a HashMap ready to use.
// The New constructor must be called to generate a new HashMap.
type HashMap struct {
	// buckets is the hash table represented by a slice.
	// The length of the slice is always equal to its capacity, that is, cap(buckets) equals len(buckets).
	buckets []*node

	// cap is the number of buckets.
	// len is the number of entries.
	cap, len int

	// loadFactor is a measure of how full the hash table is allowed to get before its capacity is automatically
	//increased.
	loadFactor float64
}

// New returns a new HashMap ready to use.
// If 'cap' is less than or equal to zero, then it will be set from its default value. The same applies to 'loadFactor'.
// 'cap' and 'loadFactor' can never be changed manually.
// Time complexity: 0(1).
func New(cap int, loadFactor float64) *HashMap {
	if cap <= 0 {
		cap = DefaultCapacity
	}
	if loadFactor <= 0 {
		loadFactor = DefaultLoadFactor
	}
	return &HashMap{
		buckets:    make([]*node, cap, cap),
		cap:        cap,
		len:        0,
		loadFactor: loadFactor,
	}
}

// NewByMap returns a new HashMap with the values stored in the map.
// Time complexity: O(n), where n is the length of the map.
func NewByMap(values map[coll.Hashable]interface{}, cap int, loadFactor float64) *HashMap {
	hm := New(cap, loadFactor)
	for k, v := range values {
		hm.Push(k, v)
	}
	return hm
}

// Clone returns a new cloned HashMap.
// Time complexity: O(c + e), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) Clone() *HashMap {
	clone := New(hm.cap, hm.loadFactor)
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			clone.Push(n.key, n.value)
		}
	}
	return clone
}

// Do gets a value and performs all the procedures, then repeats this with the rest of the values.
// The choice of values is not predictable.
// The hash map retains its original state.
// Time complexity: O(c + e), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) Do(procedures ...func(v interface{})) {
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			for _, procedure := range procedures {
				procedure(n.value)
			}
		}
	}
}

// Get returns the paired value to 'key'.
// If the hash map is empty or 'key' is not found, then returns nil and false.
// Time complexity: θ(1), assuming the hash function disperses the values properly among the buckets.
func (hm *HashMap) Get(key coll.Hashable) (v interface{}, ok bool) {
	if key == nil || hm.IsEmpty() {
		return nil, false
	}
	hash := hm.hash(key.Hash())
	if found := hm.buckets[hash].search(key); found == nil {
		return nil, false
	} else {
		return found.value, true
	}
}

// hash returns an integer representing the index where a new value will be stored in the buckets.
// Time complexity: O(1).
func (hm *HashMap) hash(hashCode int) int {
	return hashCode % hm.cap
}

// IsEmpty returns true if the hash map has no values.
// Time complexity: O(1).
func (hm *HashMap) IsEmpty() bool {
	return hm.len == 0
}

// Len returns the current length (number of entries) of the hash map.
// Time complexity: O(1).
func (hm *HashMap) Len() int {
	return hm.len
}

// Map returns a new map with the values stored in the hash map.
// The hash map retains its original state.
// Time complexity: O(c + e), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) Map() map[coll.Hashable]interface{} {
	m := make(map[coll.Hashable]interface{})
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			m[n.key] = n.value
		}
	}
	return m
}

// Push inserts the key-value pair and returns true.
// If 'key' already exists, then updates the matched value and returns true.
// The nil key is not allowed, if 'key' is nil, then returns false and does nothing.
// Time complexity: θ(1), assuming the hash function disperses the values properly among the buckets.
func (hm *HashMap) Push(key coll.Hashable, v interface{}) bool {
	if key == nil {
		return false
	}
	if lenF, capF := float64(hm.len), float64(hm.cap); lenF > capF*hm.loadFactor {
		hm.reHashing()
	}
	hashCode := key.Hash()
	hash := hm.hash(hashCode)
	if hm.buckets[hash] != nil {
		if found := hm.buckets[hash].search(key); found != nil {
			found.hashCode = hashCode
			found.key = key
			found.value = v
		} else {
			newNode := &node{
				hashCode: hashCode,
				key:      key,
				value:    v,
				next:     hm.buckets[hash],
			}
			hm.buckets[hash] = newNode
			hm.len++
		}
	} else {
		newNode := &node{
			hashCode: hashCode,
			key:      key,
			value:    v,
			next:     nil,
		}
		hm.buckets[hash] = newNode
		hm.len++
	}
	return true
}

// reHashing will double the buckets capacity and reinsert the values.
// Time complexity: O(c + e*θ(1)), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) reHashing() {
	old := hm.buckets
	hm.buckets = make([]*node, hm.cap*2, hm.cap*2)
	hm.cap *= 2
	hm.len = 0
	for _, n := range old {
		for n != nil {
			next := n.next
			hm.Push(n.key, n.value)
			n.clear()
			n = next
		}
	}
}

// Remove removes the key-value pair that matches the 'key' parameter.
// Time complexity: θ(1), assuming the hash function disperses the values properly among the buckets.
func (hm *HashMap) Remove(key coll.Hashable) (v interface{}, ok bool) {
	if key == nil {
		return nil, false
	}
	hash := hm.hash(key.Hash())
	n := hm.buckets[hash]
	if n == nil {
		return nil, false
	}
	if n.key.Equals(key) {
		v := n.value
		hm.buckets[hash] = n.next
		n.clear()
		hm.len--
		return v, true
	}

	for n.next != nil && !n.next.key.Equals(key) {
		n = n.next
	}
	if n.next != nil {
		v := n.next.value
		toRemove := n.next
		n.next = toRemove.next
		toRemove.clear()
		hm.len--
		return v, true
	}
	return nil, false
}

// RemoveAll sets the properties of the hash map to its zero values.
// Time complexity: O(1).
func (hm *HashMap) RemoveAll() {
	hm.buckets, hm.cap, hm.len, hm.loadFactor = nil, 0, 0, 0
}

// Search returns the key of the first match of the value 'v'.
// If the value 'v' does not belong to the hash map, then returns nil.
// Time complexity: O(c + e), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) Search(v interface{}) coll.Hashable {
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			if n.value == v {
				return n.key
			}
		}
	}
	return nil
}

// SearchByComparator returns the key of the first match of the value 'v'.
// If the value 'v' does not belong to the hash map, then returns nil.
// The comparison between values is defined by the parameter 'equals'.
// The function 'equals' must return true if 'v1' equals 'v2'.
// Time complexity: O(c + e), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) SearchByComparator(v interface{}, equals func(v1, v2 interface{}) bool) coll.Hashable {
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			if equals(n.value, v) {
				return n.key
			}
		}
	}
	return nil
}

// String returns a representation of the hash map as a string.
// HashMap implements the fmt.Stringer interface.
// Time complexity: O(c + e), where c is the capacity of the hash map and e its number of entries.
func (hm *HashMap) String() string {
	if hm.IsEmpty() {
		return "[]"
	}
	str := "["
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			str += fmt.Sprintf("%v:%v ", n.key, n.value)
		}
	}
	return str[:len(str)-1] + "]"
}
