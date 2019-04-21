package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"
	"strings"

	"github.com/beito123/binary"
	lvldb "github.com/beito123/goleveldb/leveldb"
	"github.com/beito123/goleveldb/leveldb/filter"
	"github.com/beito123/goleveldb/leveldb/opt"
	"github.com/beito123/level"
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
		chunks:   make(map[int]*Chunk),
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
		chunks:   make(map[int]*Chunk),
	}, nil
}

// NewBlockState returns new BlockState
func NewBlockState(name string, value int) *BlockState {
	return &BlockState{
		name:  strings.ToLower(name),
		value: value,
	}
}

// BlockState is a block information
type BlockState struct {
	name  string
	value int
}

// Name returns block name
func (block *BlockState) Name() string {
	return block.name
}

// Value returns block value
func (block *BlockState) Value() int {
	return block.value
}

// Equal returns whether block is equal b
func (block *BlockState) Equal(b *BlockState) bool {
	return block.name == b.name && block.value == b.value
}

// LevelDB is a level format for mcbe
type LevelDB struct {
	Database *lvldb.DB

	Format    ChunkFormat
	Dimension level.Dimension

	chunks map[int]*Chunk
}

func (lvl *LevelDB) at(x, y int) int {
	return y<<16 | x
}

// LoadChunk loads a chunk.
// If create is true, generates a chunk.
func (lvl *LevelDB) LoadChunk(x, y int) error {
	if lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: already loaded the chunk")
	}

	chunk, err := lvl.Format.Read(x, y, lvl.Dimension, lvl.Database)
	if err != nil {
		return err
	}

	lvl.chunks[lvl.at(x, y)] = chunk

	return nil
}

// UnloadChunk unloads a chunk.
func (lvl *LevelDB) UnloadChunk(x, y int) error {
	return nil
}

// GenerateChunk generates a chunk
func (lvl *LevelDB) GenerateChunk(x, y int) error {
	return nil
}

// HasGeneratedChunk returns whether the chunk is generaged
func (lvl *LevelDB) HasGeneratedChunk(x, y int) bool {
	return false
}

// IsLoadedChunk returns weather a chunk is loaded.
func (lvl *LevelDB) IsLoadedChunk(x, y int) bool {
	_, ok := lvl.chunks[lvl.at(x, y)]

	return ok
}

// SaveChunk saves a chunk.
func (lvl *LevelDB) SaveChunk(x, y int) error {
	return nil
}

// SaveChunks saves all chunks.
func (lvl *LevelDB) SaveChunks() error {
	return nil
}

// Chunk returns a loaded chunk.
func (lvl *LevelDB) Chunk(x, y int) (level.Chunk, bool) {
	chunk, ok := lvl.chunks[lvl.at(x, y)]

	return chunk, ok
}

// LoadedChunks returns loaded chunks.
func (lvl *LevelDB) LoadedChunks() []level.Chunk {
	result := make([]level.Chunk, len(lvl.chunks))

	count := 0
	for _, chunk := range lvl.chunks {
		result[count] = chunk
		count++
	}

	return result
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

// ChunkFormatV120 is a chunk format after v1.2.0
type ChunkFormatV120 struct {
	RuntimeIDList map[int]*BlockState
}

func (format *ChunkFormatV120) Read(x, y int, dimension level.Dimension, db *lvldb.DB) (*Chunk, error) {
	chunk := NewChunk(x, y)
	chunk.DefaultBlock = NewBlockState("minecraft:air", 0)

	stateKey := getChunkKey(x, y, dimension, TagFinalizedState, 0)

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

	// Load subchunks
	for i := 0; i < 16; i++ {
		key := getChunkKey(x, y, dimension, TagSubChunkPrefix, i)

		ok, err := db.Has(key, nil)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		val, err := db.Get(key, nil)
		if err != nil {
			return nil, err
		}

		sub, err := format.ReadSubChunk(byte(i), val)
		if err != nil {
			return nil, err
		}

		//fmt.Printf("%#v", sub)

		chunk.subChunks[i] = sub
	}

	return chunk, nil
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
			return nil, err
		}
	default: // 0, 2, 3, 4, 5, 6, 7 // 1.2
	}

	return sub, nil
}

func getChunkKey(x int, y int, dimension level.Dimension, tag byte, sid int) []byte {
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
