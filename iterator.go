// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package collection

const (
	ErrorIteratorNext    = "iterator: use HasNext method before Next"
	ErrorIteratorHasNext = "iterator: HasNext method returned false on last call"
	ErrorIteratorRemove  = "iterator: use Next method before Remove"
)

type Iterator interface {
	ForEach(action func(v *interface{}))
	HasNext() bool
	Next() (interface{}, error)
	Remove() error
}
