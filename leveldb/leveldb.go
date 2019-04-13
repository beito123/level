package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	lvldb "github.com/beito123/goleveldb/leveldb"
	"github.com/beito123/goleveldb/leveldb/filter"
	"github.com/beito123/goleveldb/leveldb/opt"
	"github.com/beito123/level"
	"github.com/beito123/level/block"
)

var DefaultOptions = &opt.Options{
	Filter:      filter.NewBloomFilter(10),
	WriteBuffer: 4 * 1024 * 1024, // 4mb
}

func New(path string) (*LevelDB, error) {
	return NewWithOptions(path, DefaultOptions)
}

func NewWithOptions(path string, options *opt.Options) (*LevelDB, error) {
	db, err := lvldb.RecoverFile(path, options)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		Database: db,
	}, nil
}

func Load(path string) (*LevelDB, error) {
	return LoadWithOptions(path, DefaultOptions)
}

func LoadWithOptions(path string, options *opt.Options) (*LevelDB, error) {
	db, err := lvldb.OpenFile(path, options)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		Database: db,
	}, nil
}

// LevelDB is a level format for mcbe
type LevelDB struct {
	Database *lvldb.DB
}

// LoadChunk loads a chunk.
// If create is true, generates a chunk.
func (lvl *LevelDB) LoadChunk(x, y int, create bool) bool {
	return false
}

// UnloadChunk unloads a chunk.
// If safe is false, unloads even when players are there.
func (lvl *LevelDB) UnloadChunk(x, y int, safe bool) bool {
	return false
}

// IsLoadedChunk returns weather a chunk is loaded.
func (lvl *LevelDB) IsLoadedChunk(x, y int) bool {
	return false
}

// SaveChunk saves a chunk.
func (lvl *LevelDB) SaveChunk(x, y int) bool {
	return false
}

// SaveChunks saves all chunks.
func (lvl *LevelDB) SaveChunks() {

}

// Chunk returns a chunk.
// If it's not loaded, loads the chunks.
// If create is true, generates a chunk.
func (lvl *LevelDB) Chunk(x, y int, create bool) level.Chunk {
	return nil
}

// Chunks retuns loaded chunks.
func (lvl *LevelDB) Chunks() []level.Chunk {
	return nil
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
	Read(x, y, dimension level.Dimension, db *lvldb.DB) (*Chunk, error)
	Write(chunk *Chunk, db *lvldb.DB)
}

// ChunkFormatV120 is a chunk format after v1.2.0
type ChunkFormatV120 struct {
	RuntimeIDList map[int]block.BlockData
}

func (format *ChunkFormatV120) Read(x, y int, dimension level.Dimension, db *lvldb.DB) (*Chunk, error) {
	chunk := NewChunk(x, y)

	// Load subchunks
	for i := 0; i < 16; i++ {
		key := getChunkKey(x, y, dimension, TagSubChunkPrefix, i)
		val, err := db.Get(key, nil)
		if err != nil {
			return nil, err
		}

		sub, err := format.ReadSubChunk(byte(i), val)
		if err != nil {
			return nil, err
		}

		chunk.subChunks[i] = sub
	}

	return nil, nil
}

func (format *ChunkFormatV120) ReadSubChunk(y byte, b []byte) (sub *SubChunk, err error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("level.leveldb: enough bytes")
	}

	ver := b[0]

	switch ver {
	case 1: // 1.2.13 only
	case 8: // after 1.3
		subFormat := &SubChunkFormatV130{
			RuntimeIDList: format.RuntimeIDList,
		}
		sub, err = subFormat.Read(y, b)
		if err != nil {
			return nil, nil
		}
	default: // 0, 2, 3, 4, 5, 6, 7 // 1.2
	}

	return sub, nil
}

type ChunkFormatV090 struct { // TODO?
}

func getChunkKey(x int, y int, dimension level.Dimension, tag byte, sid int) []byte {
	base := []byte{
		byte(x >> 24),
		byte(x >> 16),
		byte(x >> 8),
		byte(x),
		byte(y >> 24),
		byte(y >> 16),
		byte(y >> 8),
		byte(y),
	}

	switch {
	case dimension != level.OverWorld && tag == TagSubChunkPrefix:
		return []byte{
			base[0],
			base[1],
			base[2],
			base[3],
			base[4],
			base[5],
			base[6],
			base[7],
			byte(dimension),
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
			byte(dimension),
			tag,
		}
	case tag == TagSubChunkPrefix:
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

	return base
}
