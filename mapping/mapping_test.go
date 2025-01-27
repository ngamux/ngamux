package mapping

import (
	"testing"

	"github.com/golang-must/must"
)

func TestNew(t *testing.T) {
	actual := New[string, int]()
	expectedSliceLen := 0
	expectedMapLen := 0
	must.Equal(t, expectedSliceLen, len(actual.s))
	must.Equal(t, expectedMapLen, len(actual.m))
}

func TestSet(t *testing.T) {
	testCases := []struct {
		desc             string
		actual           Mapping[string, int]
		expectedSliceLen int
		expectedMapLen   int
	}{
		{
			desc: "positive",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("b", 2)
				mapp.Set("c", 3)
				return mapp
			}(),
			expectedSliceLen: 3,
			expectedMapLen:   0,
		},
		{
			desc: "duplicate key",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("a", 2)
				mapp.Set("a", 3)
				return mapp
			}(),
			expectedSliceLen: 1,
			expectedMapLen:   0,
		},
		{
			desc: "to map",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("b", 2)
				mapp.Set("c", 3)
				mapp.Set("d", 4)
				mapp.Set("e", 5)
				mapp.Set("f", 6)
				mapp.Set("g", 7)
				mapp.Set("h", 8)
				mapp.Set("i", 9)
				mapp.Set("j", 10)
				mapp.Set("k", 11)
				return mapp
			}(),
			expectedSliceLen: 0,
			expectedMapLen:   11,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			must.Equal(t, tC.expectedSliceLen, len(tC.actual.s))
			must.Equal(t, tC.expectedMapLen, len(tC.actual.m))
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		desc          string
		actual        func() (int, bool)
		expectedOk    bool
		expectedValue int
	}{
		{
			desc: "positive",
			actual: func() (int, bool) {
				mapp := New[string, int]()
				mapp.s = []mappingEntry[string, int]{
					{"a", 1},
					{"b", 2},
					{"c", 3},
				}
				return mapp.Get("a")
			},
			expectedOk:    true,
			expectedValue: 1,
		},
		{
			desc: "not exists",
			actual: func() (int, bool) {
				mapp := New[string, int]()
				return mapp.Get("a")
			},
			expectedOk:    false,
			expectedValue: 0,
		},
		{
			desc: "duplicate key",
			actual: func() (int, bool) {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("a", 2)
				mapp.Set("a", 3)
				return mapp.Get("a")
			},
			expectedOk:    true,
			expectedValue: 3,
		},
		{
			desc: "to map",
			actual: func() (int, bool) {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("b", 2)
				mapp.Set("c", 3)
				mapp.Set("d", 4)
				mapp.Set("e", 5)
				mapp.Set("f", 6)
				mapp.Set("g", 7)
				mapp.Set("h", 8)
				mapp.Set("i", 9)
				mapp.Set("j", 10)
				mapp.Set("k", 11)
				return mapp.Get("k")
			},
			expectedOk:    true,
			expectedValue: 11,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, actualOk := tC.actual()
			must.Equal(t, tC.expectedOk, actualOk)
			must.Equal(t, tC.expectedValue, actual)
		})
	}
}

func TestEach(t *testing.T) {
	testCases := []struct {
		desc        string
		actual      Mapping[string, int]
		eachFunc    func(*int) func(string, int) bool
		expectedLen int
	}{
		{
			desc: "positive",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("b", 2)
				mapp.Set("c", 3)
				return mapp
			}(),
			eachFunc: func(actualLen *int) func(string, int) bool {
				*actualLen = 0
				return func(k string, v int) bool {
					*actualLen++
					return true
				}
			},
			expectedLen: 3,
		},
		{
			desc: "duplicate key",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("a", 2)
				mapp.Set("a", 3)
				return mapp
			}(),
			eachFunc: func(actualLen *int) func(string, int) bool {
				*actualLen = 0
				return func(k string, v int) bool {
					*actualLen++
					return true
				}
			},
			expectedLen: 1,
		},
		{
			desc: "to map",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("b", 2)
				mapp.Set("c", 3)
				mapp.Set("d", 4)
				mapp.Set("e", 5)
				mapp.Set("f", 6)
				mapp.Set("g", 7)
				mapp.Set("h", 8)
				mapp.Set("i", 9)
				mapp.Set("j", 10)
				mapp.Set("k", 11)
				return mapp
			}(),
			eachFunc: func(actualLen *int) func(string, int) bool {
				*actualLen = 0
				return func(k string, v int) bool {
					*actualLen++
					return true
				}
			},
			expectedLen: 11,
		},
		{
			desc: "with skip",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				return mapp
			}(),
			eachFunc: func(actualLen *int) func(string, int) bool {
				*actualLen = 0
				return func(k string, v int) bool {
					*actualLen++
					return false
				}
			},
			expectedLen: 1,
		},
		{
			desc: "with skip to map",
			actual: func() Mapping[string, int] {
				mapp := New[string, int]()
				mapp.Set("a", 1)
				mapp.Set("b", 2)
				mapp.Set("c", 3)
				mapp.Set("d", 4)
				mapp.Set("e", 5)
				mapp.Set("f", 6)
				mapp.Set("g", 7)
				mapp.Set("h", 8)
				mapp.Set("i", 9)
				mapp.Set("j", 10)
				mapp.Set("k", 11)
				return mapp
			}(),
			eachFunc: func(actualLen *int) func(string, int) bool {
				*actualLen = 0
				return func(k string, v int) bool {
					*actualLen++
					return false
				}
			},
			expectedLen: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actualLen := 0
			tC.actual.Each(tC.eachFunc(&actualLen))
			must.Equal(t, tC.expectedLen, actualLen)
		})
	}
}
