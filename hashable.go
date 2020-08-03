// Copyright 2020 maguerrido <mauricio.aguerrido@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package collection

// Hashable defines a data type capable of generating its own hash code and also capable of being compared with another
//hashable value.
type Hashable interface {
	// Equals returns true if 'v' equals the caller value.
	Equals(v Hashable) bool

	// Hash returns the hash code of the caller value.
	Hash() int
}
