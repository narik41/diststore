package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, any) error
}

type GOBDecoder struct{}

func NewGOBDecoder() Decoder {
	return &GOBDecoder{}
}

func (dec GOBDecoder) Decode(r io.Reader, obj any) error {
	return gob.NewDecoder(r).Decode(obj)
}
