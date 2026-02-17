package json

import (
	stdjson "encoding/json"
	"errors"
	"testing"

	"github.com/golang-must/must"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMarshalUnmarshal_Default(t *testing.T) {
	m := must.New(t)
	p := Person{Name: "Alice", Age: 30}
	got, err := Marshal(p)
	m.Nil(err)
	want, err := stdjson.Marshal(p)
	m.Nil(err)
	m.Equal(string(got), string(want))
	var out Person
	m.Nil(Unmarshal(got, &out))
	m.Equal(out, p)
}

func TestConfigure_OverrideMarshalUnmarshal(t *testing.T) {
	m := must.New(t)
	origMarshal := globalMarshaler
	origUnmarshal := globalUnmarshaler
	defer func() {
		Configure(origMarshal, origUnmarshal)
	}()
	mockMarshal := func(v any) ([]byte, error) {
		return []byte("__mocked__"), nil
	}
	mockUnmarshal := func(data []byte, v any) error {
		switch out := v.(type) {
		case *Person:
			out.Name = "Injected"
			out.Age = 99
			return nil
		default:
			return errors.New("unsupported type")
		}
	}
	Configure(mockMarshal, mockUnmarshal)
	b, err := Marshal(Person{Name: "Bob", Age: 1})
	m.Nil(err)
	m.Equal(string(b), "__mocked__")
	var p Person
	m.Nil(Unmarshal([]byte("irrelevant"), &p))
	m.Equal(p, Person{Name: "Injected", Age: 99})
}
