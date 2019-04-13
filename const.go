package level

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

type Dimension int

func (d Dimension) ToMCBE() byte {
	switch d {
	case OverWorld:
	}

	return 0
}

const (
	OverWorld Dimension = iota
	Nether
	TheEnd
	Unknown
)
