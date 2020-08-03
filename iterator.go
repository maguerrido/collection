// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package collection

const (
	// ErrorIteratorNext will be returned when Next is called without first calling HasNext.
	ErrorIteratorNext = "iterator: use HasNext method before Next"

	// ErrorIteratorHasNext will be returned when Next is called with a false return of HasNext.
	ErrorIteratorHasNext = "iterator: HasNext method returned false on last call"

	// ErrorIteratorRemove will be returned when Remove is called without first calling Next.
	ErrorIteratorRemove = "iterator: use Next method before Remove"

	// ErrorIteratorRemoveNotSupported will be returned if the abstract data type does not support this method.
	ErrorIteratorRemoveNotSupported = "iterator: Remove method not supported"
)

// Iterator defines a data type capable of traversing an entire collection of data.
type Iterator interface {
	// ForEach modifies all the values stored in the collection.
	ForEach(action func(v *interface{}))

	// HasNext returns true if the iterator has not yet finished browsing the entire collection.
	HasNext() bool

	// Next return the next value in the collection.
	Next() (interface{}, error)

	// Remove removes the value pointed by the iterator (the last Next call).
	Remove() error
}
