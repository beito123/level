package level

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

// Dimension is a kind of world type
type Dimension int

const (
	// OverWorld is a first dimension at new game
	OverWorld Dimension = iota

	// Nether is a hell world
	Nether

	// TheEnd is a dimension with floating islands
	TheEnd

	// Unknown is for unknown dimension
	// It' not use generally
	Unknown
)
