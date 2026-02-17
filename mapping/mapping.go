// Package mapping provides a small, generic key->value container optimized
// for a small number of entries. For up to maxSlice entries it stores data
// in a small slice (to avoid map allocations). Once the number of entries
// grows beyond that threshold it upgrades to a map for O(1) lookups. This
// optimization balances allocation cost for small maps and lookup speed for
// larger collections.
package mapping

const maxSlice = 10

type mappingEntry[K comparable, V any] struct {
	k K
	v V
}

// Mapping is a generic container that holds key/value pairs. It starts
// using a small slice for storage and transparently migrates to a map
// when the number of entries exceeds maxSlice.
type Mapping[K comparable, V any] struct {
	s []mappingEntry[K, V]
	m map[K]V
}

// New constructs an empty Mapping.
func New[K comparable, V any]() Mapping[K, V] {
	return Mapping[K, V]{}
}

// Set inserts or updates the value for a given key. While the container is
// small it stores entries in the slice; once capacity surpasses maxSlice it
// migrates existing entries into a map and proceeds using the map.
func (mapp *Mapping[K, V]) Set(k K, v V) {
	if mapp.m == nil && len(mapp.s) < maxSlice {
		found := -1
		for i, e := range mapp.s {
			if e.k == k {
				found = i
				break
			}
		}
		if found <= -1 {
			mapp.s = append(mapp.s, mappingEntry[K, V]{k, v})
			return
		}
		mapp.s[found] = mappingEntry[K, V]{k, v}
		return
	}

	if mapp.m == nil {
		mapp.m = map[K]V{}
		for _, e := range mapp.s {
			mapp.m[e.k] = e.v
		}
		mapp.s = nil
	}

	mapp.m[k] = v
}

// Get retrieves the value for the provided key. The boolean return value
// indicates whether the key was present.
func (mapp *Mapping[K, V]) Get(k K) (v V, ok bool) {
	if mapp.m != nil {
		v, ok = mapp.m[k]
		return v, ok
	}

	for _, e := range mapp.s {
		if e.k == k {
			return e.v, true
		}
	}
	return v, ok
}

// Each iterates over all entries in the mapping and calls fn for each key
// and value. If fn returns false, iteration stops early. The iteration
// order is undefined for map-backed storage and insertion order for
// slice-backed storage.
func (mapp *Mapping[K, V]) Each(fn func(k K, v V) bool) {
	if mapp.m != nil {
		for k, v := range mapp.m {
			if !fn(k, v) {
				return
			}
		}
	} else {
		for _, e := range mapp.s {
			if !fn(e.k, e.v) {
				return
			}
		}
	}

}
