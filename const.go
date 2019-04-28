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

/*
// HeightMapType returns type of heightmap
type HeightMapType int

const (
	// MotionBlocking contains blocks block motion and a fluid
	MotionBlocking HeightMapType = iota

	// MotionBlockingNoLeaves contains blocks block motion, a fluid and leaves
	MotionBlockingNoLeaves

	// OceanFloor contains non air and soild block
	OceanFloor

	// OceanFloorWorldGeneration contains non air and fluid block. For world generation
	OceanFloorWorldGeneration

	// WorldSurface contains non air block
	WorldSurface

	// WorldSurfaceWorldGeneration contains non air block. For world generation
	WorldSurfaceWorldGeneration
)

// Height returns the height of the highest block at chunk coordinate
// If kind is not supported for a format, returns false for ok
Height(x, y int, kind HeightMapType) (height uint16, ok bool)*/
