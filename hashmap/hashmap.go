// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package hashmap

import (
	coll "github.com/maguerrido/collection"
)

const (
	DefaultCapacity   = 16
	DefaultLoadFactor = 0.75
)

type node struct {
	hashCode int
	key      coll.Hashable
	value    interface{}
	next     *node
}

func (n *node) clear() {
	n.hashCode, n.key, n.value, n.next = 0, nil, nil, nil
}

func (n *node) search(k coll.Hashable) *node {
	for i := n; i != nil; i = i.next {
		if i.key.Equals(k) {
			return i
		}
	}
	return nil
}

type HashMap struct {
	buckets    []*node
	cap, len   int
	loadFactor float64
}

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

func NewByMap(values map[coll.Hashable]interface{}, cap int, loadFactor float64) *HashMap {
	hm := New(cap, loadFactor)
	for k, v := range values {
		hm.Push(k, v)
	}
	return hm
}

func (hm *HashMap) Clone() *HashMap {
	clone := New(hm.cap, hm.loadFactor)
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			clone.Push(n.key, n.value)
		}
	}
	return clone
}

func (hm *HashMap) Get(key coll.Hashable) (v interface{}, ok bool) {
	if key == nil {
		return nil, false
	}
	hash := hm.hash(key.Hash())
	if found := hm.buckets[hash].search(key); found == nil {
		return nil, false
	} else {
		return found.value, true
	}
}

func (hm *HashMap) hash(hashCode int) int {
	return hashCode % hm.cap
}

func (hm *HashMap) IsEmpty() bool {
	return hm.len == 0
}

func (hm *HashMap) Map() map[coll.Hashable]interface{} {
	m := make(map[coll.Hashable]interface{})
	for _, n := range hm.buckets {
		for ; n != nil; n = n.next {
			m[n.key] = n.value
		}
	}
	return m
}

func (hm *HashMap) Len() int {
	return hm.len
}

func (hm *HashMap) Push(key coll.Hashable, v interface{}) bool {
	if key == nil {
		return false
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

func (hm *HashMap) RemoveAll() {
	hm.buckets, hm.cap, hm.len, hm.loadFactor = nil, 0, 0, 0
}
