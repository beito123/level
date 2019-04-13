package mcbe

import (
	"github.com/beito123/level"
)

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

type DimensionID int

const (
	OverWorld DimensionID = iota
	Nether                // Required verification
	TheEnd                // Required verification

	Unknown DimensionID = -1
)

func ToDimensionID(d level.Dimension) DimensionID {
	switch d {
	case level.OverWorld:
		return OverWorld
	case level.Nether:
		return Nether
	case level.TheEnd:
		return TheEnd
	}

	return Unknown
}

func FromDimensionID(id DimensionID) level.Dimension {
	switch id {
	case OverWorld:
		return level.OverWorld
	case Nether:
		return level.Nether
	case TheEnd:
		return level.TheEnd
	}

	return level.Unknown
}
