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

func (hm *HashMap) hash(hashCode int) int {
	return hashCode % hm.cap
}

func (hm *HashMap) IsEmpty() bool {
	return hm.len == 0
}

func (hm *HashMap) Len() int {
	return hm.len
}

func (hm *HashMap) RemoveAll() {
	hm.buckets, hm.cap, hm.len, hm.loadFactor = nil, 0, 0, 0
}
