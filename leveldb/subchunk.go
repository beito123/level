package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import "fmt"

// BlockStorageSize is a size of BlockStorage
const BlockStorageSize = 16 * 16 * 16

// NewBlockStorage returns new BlockStorage
func NewBlockStorage() *BlockStorage {
	return &BlockStorage{
		Blocks: make([]uint16, BlockStorageSize),
	}
}

// BlockStorage is a storage contains blocks of a subchunk
type BlockStorage struct {
	Palettes []*BlockState
	Blocks   []uint16
}

// At returns a index for Blocks at blockstorage coordinates
func (BlockStorage) At(x, y, z int) int {
	return x<<8 | z<<4 | y
}

// Vaild vailds blockstorage coordinates
func (BlockStorage) Vaild(x, y, z int) error {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return fmt.Errorf("level.leveldb: invail coordinate")
	}

	return nil
}

// GetBlock returns the BlockState at blockstorage coordinates
func (storage *BlockStorage) GetBlock(x, y, z int) (*BlockState, error) {
	err := storage.Vaild(x, y, z)
	if err != nil {
		return nil, err
	}

	index := storage.At(x, y, z)

	if index >= len(storage.Blocks) {
		return nil, fmt.Errorf("level.leveldb: uninitialized BlockStorage")
	}

	id := storage.Blocks[index]

	if int(id) >= len(storage.Palettes) {
		return nil, fmt.Errorf("level.leveldb: couldn't find a palette for the block")
	}

	return storage.Palettes[id], nil
}

// NewSubChunk returns new SubChunk
func NewSubChunk(y byte) *SubChunk {
	return &SubChunk{
		Y: y,
	}
}

// SubChunk is a 16x16x16 blocks segment for a chunk
type SubChunk struct {
	Y byte

	Storages []*BlockStorage
}

// GetBlockStorage returns BlockStorage which subchunk contained with index
func (sub *SubChunk) GetBlockStorage(index int) (*BlockStorage, bool) {
	if index >= len(sub.Storages) || index < 0 {
		return nil, false
	}

	return sub.Storages[index], true
}

// AtBlock returns BlockState at the subchunk coordinates
func (sub *SubChunk) AtBlock(x, y, z, index int) (*BlockState, error) {
	storage, ok := sub.GetBlockStorage(index)
	if !ok {
		return nil, fmt.Errorf("level.leveldb: invaild storage index")
	}

	return storage.GetBlock(x, y, z)
}

// SubChunkFormat is a formatter for subchunk
type SubChunkFormat interface {
	Read(y byte, b []byte) (*SubChunk, error)
}
