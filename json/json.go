package json

import "encoding/json"

var globalMarshaler = json.Marshal
var globalUnmarshaler = json.Unmarshal

func Configure(marshaler func(any) ([]byte, error), unmarshaler func([]byte, any) error) func() {
	if marshaler != nil {
		globalMarshaler = marshaler
	}

	if unmarshaler != nil {
		globalUnmarshaler = unmarshaler
	}

	return func() {}
}

func Marshal(v any) ([]byte, error) {
	return globalMarshaler(v)
}

func Unmarshal(data []byte, v any) error {
	return globalUnmarshaler(data, v)
}
