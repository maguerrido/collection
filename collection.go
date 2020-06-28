// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package collection

type Hashable interface {
	Equals(v Hashable) bool
	Hash() int
}
