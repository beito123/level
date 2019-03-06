package binary

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

// TriadSize is the Triad length(bytes)
const TriadSize = 3

const (

	// MinTriad is the minimum value of Triad
	MinTriad = 0

	// MaxTriad is the maximum value of Triad
	MaxTriad = 16777216
)

// Triad is a type of 3bytes
type Triad uint32

// ReadTriad reads Triad value
func ReadTriad(v []byte) Triad {
	return Triad(v[0])<<16 | Triad(v[1])<<8 | Triad(v[2])
}

// WriteTriad writes Triad value
func WriteTriad(v Triad) []byte {
	return []byte{
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

// ReadLTriad reads Triad value as LittleEndian
func ReadLTriad(v []byte) Triad {
	return Triad(v[0]) | Triad(v[1])<<8 | Triad(v[2])<<16
}

// WriteTriad writes Triad value as LittleEndian
func WriteLTriad(v Triad) []byte {
	return []byte{
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
	}
}

// ReadETriad reads Triad value with error
func ReadETriad(v []byte) (Triad, error) {
	if len(v) < TriadSize {
		return 0, nil
	}

	return ReadTriad(v), nil
}

// ReadELTriad reads Triad value as LittleEndian with error
func ReadELTriad(v []byte) (Triad, error) {
	if len(v) < TriadSize {
		return 0, nil
	}

	return ReadLTriad(v), nil
}
