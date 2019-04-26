package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import "fmt"
import "math"

// GetStorageTypeFromSize returns 
func GetStorageTypeFromSize(size uint) StorageType {
	size--
	size |= (size >> 1)
	size |= (size >> 2)
	size |= (size >> 4)
	size |= (size >> 8)
	size |= (size >> 16)
	size++

	return StorageType(math.Log2(float64(size)))
}

// StorageType is a type of BlockStorage
type StorageType int

// BitsPerBlock retunrs bits per a block for BlockStorage
func (t StorageType) BitsPerBlock() int {
	return int(t)
}

// PaletteSize returns a size of palette for StorageType
func (t StorageType) PaletteSize() int {
	return 1 << uint(t)
}

const (
	TypePalette1 StorageType = 1
	TypePalette2 StorageType = 2
	TypePalette3 StorageType = 3
	TypePalette4 StorageType = 4
	TypePalette5 StorageType = 5
	TypePalette6 StorageType = 6
	TypePalette8 StorageType = 8
	TypePalette16 StorageType = 16
)

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
	Palettes []*RawBlockState
	Blocks   []uint16
}

// At returns a index for Blocks at blockstorage coordinates
func (BlockStorage) At(x, y, z int) int {
	return x<<8 | z<<4 | y
}

// Vaild vailds blockstorage coordinates
func (BlockStorage) Vaild(x, y, z int) error {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return fmt.Errorf("level.leveldb: invaild block storage coordinate")
	}

	return nil
}

// GetBlock returns the BlockState at blockstorage coordinates
func (storage *BlockStorage) GetBlock(x, y, z int) (*RawBlockState, error) {
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

// SetBlock set the BlockState at blockstorage coordinates
func (storage *BlockStorage) SetBlock(x, y, z int, bs *RawBlockState) error {
	if len(storage.Palettes) > TypePalette16.PaletteSize() {
		return fmt.Errorf("level.leveldb: unsupported palette size > %d", TypePalette16.PaletteSize())
	}

	for i, v := range storage.Palettes {
		if v.Equal(bs) {
			storage.Blocks[storage.At(x, y, z)] = uint16(i)
			return nil
		}
	}

	storage.Palettes = append(storage.Palettes, bs)
	storage.Blocks[storage.At(x, y, z)] = uint16(len(storage.Palettes)-1)

	return nil
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

// GetBlock returns BlockState at the subchunk coordinates
func (sub *SubChunk) GetBlock(x, y, z, index int) (*RawBlockState, error) {
	storage, ok := sub.GetBlockStorage(index)
	if !ok {
		return nil, fmt.Errorf("level.leveldb: invaild storage index")
	}

	return storage.GetBlock(x, y, z)
}

// SetBlock returns BlockState at the subchunk coordinates
func (sub *SubChunk) SetBlock(x, y, z, index int, bs *RawBlockState) error {
	storage, ok := sub.GetBlockStorage(index)
	if !ok {
		return fmt.Errorf("level.leveldb: invaild storage index")
	}

	return storage.SetBlock(x, y, z, bs)
}

// SubChunkFormat is a formatter for subchunk
type SubChunkFormat interface {
	Read(y byte, b []byte) (*SubChunk, error)
}
