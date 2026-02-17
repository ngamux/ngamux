// Package json provides a small wrapper around Go's encoding/json
// package to allow the application to override the global marshalling
// and unmarshalling functions used throughout the Ngamux project.
//
// It exposes Configure to replace the underlying functions and thin
// helpers Marshal and Unmarshal that delegate to the configured
// implementations. This is useful when a consumer wants to use a
// different JSON implementation (for example, a faster or custom
// serializer) without changing callsites across the codebase.
package json

import "encoding/json"

// globalMarshaler is the function used by Marshal. It defaults to
// encoding/json.Marshal but can be replaced with Configure.
var globalMarshaler = json.Marshal

// globalUnmarshaler is the function used by Unmarshal. It defaults to
// encoding/json.Unmarshal but can be replaced with Configure.
var globalUnmarshaler = json.Unmarshal

// Configure replaces the global marshaler and/or unmarshaler used by
// this package. Pass nil for either parameter to leave that function
// unchanged. The returned function is a no-op cleanup placeholder to
// encourage symmetry with other Configure helpers in the codebase.
func Configure(marshaler func(any) ([]byte, error), unmarshaler func([]byte, any) error) func() {
	if marshaler != nil {
		globalMarshaler = marshaler
	}

	if unmarshaler != nil {
		globalUnmarshaler = unmarshaler
	}

	return func() {}
}

// Marshal delegates to the configured global marshaler.
func Marshal(v any) ([]byte, error) {
	return globalMarshaler(v)
}

// Unmarshal delegates to the configured global unmarshaler.
func Unmarshal(data []byte, v any) error {
	return globalUnmarshaler(data, v)
}
