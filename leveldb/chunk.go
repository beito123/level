package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

func NewChunk(x, y int) *Chunk {
	return &Chunk{
		x:         x,
		y:         y,
		biomes:    make([]byte, 256),
		subChunks: make([]*SubChunk, 16),
	}
}

type Chunk struct {
	x         int
	y         int
	biomes    []byte
	subChunks []*SubChunk
}
