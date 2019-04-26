package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	"github.com/beito123/binary"
	lvldb "github.com/beito123/goleveldb/leveldb"
	"github.com/beito123/goleveldb/leveldb/util"
	"github.com/beito123/level"
)

// DefaultStorageIndex is the default index for StorageIndex
const DefaultStorageIndex = 0

// NewChunk returns new Chunk
func NewChunk(x, y int) *Chunk {
	return &Chunk{
		x:                   x,
		y:                   y,
		biomes:              make([]byte, 256),
		subChunks:           make([]*SubChunk, 16),
		Finalization:        NotGenerated,
		DefaultStorageIndex: DefaultStorageIndex,
	}
}

// Finalization show the status of a chunk
// It's introduced in mcpe v1.1
type Finalization int

const (
	// Unsupported is unsupported finalization by the chunk format
	Unsupported Finalization = iota

	// NotGenerated is not generated a chunk if it's set
	NotGenerated

	// NotSpawnMobs is not spawned mobs if it's set
	NotSpawnMobs

	// Generated is generated a chunk if it's set
	Generated
)

// GetFinalization returns Finalization by id
func GetFinalization(id int) (Finalization, bool) {
	switch id {
	case 0:
		return NotGenerated, true
	case 1:
		return NotSpawnMobs, true
	case 2:
		return Generated, true
	}

	return Unsupported, false
}

// Chunk is a block area which splits a world by 16x16
// It has informations of block, biomes and etc...
type Chunk struct {
	x         int
	y         int
	biomes    []byte
	subChunks []*SubChunk

	Finalization Finalization

	DefaultBlock        *RawBlockState
	DefaultStorageIndex int
}

// X returns x coordinate
func (chunk *Chunk) X() int {
	return chunk.x
}

// Y returns y coordinate
func (chunk *Chunk) Y() int {
	return chunk.y
}

// SetX set x coordinate
func (chunk *Chunk) SetX(x int) {
	chunk.x = x
}

// SetY set y coordinate
func (chunk *Chunk) SetY(y int) {
	chunk.y = y
}

// SubChunks returns sub chunks
func (chunk *Chunk) SubChunks() []*SubChunk {
	return chunk.subChunks
}

// GetSubChunk returns a sub chunk
func (chunk *Chunk) GetSubChunk(index int) (*SubChunk, bool) {
	if index >= len(chunk.subChunks) {
		return nil, false
	}

	return chunk.subChunks[index], chunk.subChunks[index] != nil
}

// AtSubChunk returns a sub chunk
func (chunk *Chunk) AtSubChunk(y int) (*SubChunk, bool) {
	return chunk.GetSubChunk(y / 16)
}

// Vaild vailds a chunk coordinates
func (chunk *Chunk) Vaild(x, y, z int) bool {
	return x >= 0 && x <= 15 && y >= 0 && y <= 256 && z >= 0 && z <= 15
}

// GetBlock gets a BlockState at a chunk coordinate
func (chunk *Chunk) GetBlock(x, y, z int) (level.BlockState, error) {
	return chunk.GetBlockAtStorage(x, y, z, chunk.DefaultStorageIndex)
}

// GetBlockAtStorage gets a BlockState at a chunk coordinate from storage of index
func (chunk *Chunk) GetBlockAtStorage(x, y, z, index int) (*RawBlockState, error) {
	if !chunk.Vaild(x, y, z) {
		return nil, fmt.Errorf("level.leveldb: invaild chunk coordinate")
	}

	sub, ok := chunk.AtSubChunk(y)
	if !ok {
		return chunk.DefaultBlock, nil // Air
	}

	return sub.GetBlock(x, y&15, z, index)
}

// SetBlock set a BlockState at chunk coordinate
func (chunk *Chunk) SetBlock(x, y, z int, bs level.BlockState) error {
	rbs, err := FromRawBlockState(bs)
	if err != nil {
		return err
	}

	return chunk.SetBlockAtStorage(x, y, z, DefaultStorageIndex, rbs)
}

// SetBlockAtStorage set a BlockState at chunk coordinate to storage of index
func (chunk *Chunk) SetBlockAtStorage(x, y, z, index int, bs *RawBlockState) error {
	if !chunk.Vaild(x, y, z) {
		return fmt.Errorf("level.leveldb: invaild chunk coordinate")
	}

	sub, ok := chunk.AtSubChunk(y)
	if !ok {
		sub = NewSubChunk(byte(y/16))
	}

	return sub.SetBlock(x, y&15, z, index, bs)
}

const (
	TagData2D         = 45
	TagData2dLegacy   = 46
	TagSubChunkPrefix = 47
	TagLegacyTerrain  = 48
	TagBlockEntity    = 49
	TagEntity         = 50
	TagPendingTicks   = 51
	TagBlockExtraData = 52
	TagBiomeState     = 53
	TagFinalizedState = 54
	TagVersion        = 118
)

// ChunkFormat is a chunk format reader and writer
type ChunkFormat interface {
	// Read reads a chunk by x, y and dimension
	Read(x, y int, dimension level.Dimension, db *lvldb.DB) (*Chunk, error)
	//Write(chunk *Chunk, db *lvldb.DB)
}

// ChunkFormatV120 is a chunk format v1.2.0 or after
type ChunkFormatV120 struct {
	//RuntimeIDList map[int]*BlockState
}

func (format *ChunkFormatV120) Read(x, y int, dimension level.Dimension, db *lvldb.DB) (*Chunk, error) {
	chunk := NewChunk(x, y)
	chunk.DefaultBlock = NewRawBlockState("minecraft:air", 0)

	stateKey := format.getChunkKey(x, y, dimension, TagFinalizedState, 0)

	hasState, err := db.Has(stateKey, nil)
	if err != nil {
		return nil, err
	}

	if hasState { // after 1.1
		state, err := db.Get(stateKey, nil)
		if err != nil {
			return nil, err
		}

		if len(state) < 4 {
			return nil, fmt.Errorf("level.leveldb: invaild finalization state")
		}

		var ok bool
		chunk.Finalization, ok = GetFinalization(int(binary.ReadLInt(state)))
		if !ok {
			return nil, fmt.Errorf("level.leveldb: unknown finalization state id: %d", state)
		}
	} else {
		chunk.Finalization = Unsupported
	}

	if chunk.Finalization == NotGenerated {
		return chunk, nil
	}

	prefix := format.getChunkKey(x, y, dimension, TagSubChunkPrefix, -1)

	iter := db.NewIterator(util.BytesPrefix(prefix), nil)

	// Load subchunks
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		y := (key[len(key)-1]) & 15

		sub, err := format.ReadSubChunk(y, val)
		if err != nil {
			return nil, err
		}

		chunk.subChunks[y] = sub
	}

	return chunk, nil
}

// ReadSubChunk reads a subchunk from bytes b
func (format *ChunkFormatV120) ReadSubChunk(y byte, b []byte) (sub *SubChunk, err error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("level.leveldb: not enough bytes")
	}

	ver := b[0]

	switch ver {
	case 0, 2, 3, 4, 5, 6, 7: // v1.2 or before
		// TODO: support old format
		return nil, fmt.Errorf("level.leveldb: unsupported old subchunk format")
	case 1, 8: // Palettized format // 1.2.13 or after
		subFormat := &SubChunkFormatV1213{
			//RuntimeIDList: format.RuntimeIDList,
		}

		sub, err = subFormat.Read(y, b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("level.leveldb: unsupported subchunk version %d", ver)
	}

	return sub, nil
}

func (format *ChunkFormatV120) toDimensionID(dimension level.Dimension) int {
	switch dimension {
	case level.OverWorld:
		return 0
	case level.Nether:
		return 1
	case level.TheEnd:
		return 2
	}

	return 0
}

func (format *ChunkFormatV120) fromDimensionID(id int) level.Dimension {
	switch id {
	case 0:
		return level.OverWorld
	case 1:
		return level.Nether
	case 2:
		return level.TheEnd
	}

	return level.Unknown
}

func (format *ChunkFormatV120) getChunkKey(x int, y int, dimension level.Dimension, tag byte, sid int) []byte {
	base := []byte{
		byte(x),
		byte(x >> 8),
		byte(x >> 16),
		byte(x >> 24),
		byte(y),
		byte(y >> 8),
		byte(y >> 16),
		byte(y >> 24),
	}

	switch {
	case dimension != level.OverWorld && sid != -1:
		return []byte{
			base[0],
			base[1],
			base[2],
			base[3],
			base[4],
			base[5],
			base[6],
			base[7],
			byte(format.toDimensionID(dimension)),
			tag,
			byte(sid),
		}
	case dimension != level.OverWorld:
		return []byte{
			base[0],
			base[1],
			base[2],
			base[3],
			base[4],
			base[5],
			base[6],
			base[7],
			byte(format.toDimensionID(dimension)),
			tag,
		}
	case sid != -1:
		return []byte{
			base[0],
			base[1],
			base[2],
			base[3],
			base[4],
			base[5],
			base[6],
			base[7],
			tag,
			byte(sid),
		}
	}

	return []byte{
		base[0],
		base[1],
		base[2],
		base[3],
		base[4],
		base[5],
		base[6],
		base[7],
		tag,
	}
}
