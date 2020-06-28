// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package hashmap

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		cap        int
		loadFactor float64
		out        HashMap
	}{
		{"default_cap", 0, 0.30,
			HashMap{
				buckets:    make([]*node, DefaultCapacity, DefaultCapacity),
				cap:        DefaultCapacity,
				len:        0,
				loadFactor: 0.30,
			}},
		{"default_loadFactor", 20, 0,
			HashMap{
				buckets:    make([]*node, 20, 20),
				cap:        20,
				len:        0,
				loadFactor: DefaultLoadFactor,
			}},
		{"common", 20, 0.80,
			HashMap{
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
