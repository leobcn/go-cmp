// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package cmpopts

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/internal/value"
)

// DiscardMapZeros returns a Transformer option that discards all map entries
// with a zero value. This is useful for comparing two maps in such a way
// that entries with the zero value are ignored.
//
// DiscardMapZeros can be used in conjunction with EquateEmpty,
// but cannot be used with SortMaps.
func DiscardMapZeros() cmp.Option {
	return cmp.FilterValues(eitherHasZeroValues, cmp.Transformer("DiscardZeros", discardMapZeros))
}
func eitherHasZeroValues(x, y interface{}) bool {
	vx := reflect.ValueOf(x)
	vy := reflect.ValueOf(y)
	if x == nil || y == nil || vx.Type() != vy.Type() || vx.Kind() != reflect.Map || vx.Len()+vy.Len() == 0 {
		return false
	}
	return hasZeroValues(vx) || hasZeroValues(vy)
}
func hasZeroValues(v reflect.Value) bool {
	for _, k := range v.MapKeys() {
		if value.IsZero(v.MapIndex(k)) {
			return true
		}
	}
	return false
}
func discardMapZeros(v interface{}) interface{} {
	src := reflect.ValueOf(v)
	dst := reflect.MakeMap(src.Type())
	for _, k := range src.MapKeys() {
		if v := src.MapIndex(k); !value.IsZero(v) {
			dst.SetMapIndex(k, v)
		}
	}
	return dst.Interface()
}
