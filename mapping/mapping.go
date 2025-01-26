package mapping

const maxSlice = 10

type mappingEntry[K comparable, V any] struct {
	k K
	v V
}
type Mapping[K comparable, V any] struct {
	s []mappingEntry[K, V]
	m map[K]V
}

func New[K comparable, V any]() Mapping[K, V] {
	return Mapping[K, V]{}
}

func (mapp *Mapping[K, V]) Set(k K, v V) {
	if mapp.m == nil && len(mapp.s) < maxSlice {
		mapp.s = append(mapp.s, mappingEntry[K, V]{k, v})
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

func (mapp *Mapping[K, V]) each(fn func(k K, v V) bool) {
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
