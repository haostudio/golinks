package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/haostudio/golinks/internal/encoding"
)

// New returns a gob binary encoding/decoding instance.
func New() encoding.Binary { return &shared }

var shared dummy

type dummy struct{}

func (g *dummy) Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *dummy) Decode(b []byte, v interface{}) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}

func (g *dummy) String() string {
	return "gob"
}
