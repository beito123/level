package binary

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"github.com/beito123/binary"
)

// NewStream returns new Stream
func NewStream() *Stream {
	return NewStreamBytes([]byte{})
}

// NewStream returns new Stream with bytes
func NewStreamBytes(b []byte) *Stream {
	return &Stream{
		Stream: *binary.NewStreamBytes(b),
	}
}

// Stream is binary stream
type Stream struct {
	binary.Stream
}

// Triad sets triad got from buffer to value
func (bs *Stream) Triad() (Triad, error) {
	return ReadETriad(bs.Get(TriadSize))
}

// PutTriad puts triad from value to buffer
func (bs *Stream) PutTriad(value Triad) error {
	return bs.Put(WriteTriad(value))
}

// LTriad sets triad got from buffer as LittleEndian to value
func (bs *Stream) LTriad() (Triad, error) {
	return ReadELTriad(bs.Get(TriadSize))
}

// PutLTriad puts triad from value to buffer as LittleEndian
func (bs *Stream) PutLTriad(value Triad) error {
	return bs.Put(WriteLTriad(value))
}
