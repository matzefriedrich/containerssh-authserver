package handlers

import (
	"encoding/json"
	"io"
)

// JSONCodec defines an interface for encoding and decoding JSON data.
type JSONCodec interface {
	Decode(r io.Reader, v any) error
	Encode(w io.Writer, v any) error
}

type defaultJSONCodec struct{}

// Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
func (c *defaultJSONCodec) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

// Encode writes the JSON encoding of v to the stream.
func (c *defaultJSONCodec) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// NewJSONCodec creates a new instance of the default JSON codec.
func NewJSONCodec() JSONCodec {
	return &defaultJSONCodec{}
}
