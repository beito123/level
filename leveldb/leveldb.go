package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	lvldb "github.com/beito123/goleveldb/leveldb"
	"github.com/beito123/goleveldb/leveldb/filter"
	"github.com/beito123/goleveldb/leveldb/opt"
	"github.com/beito123/level"
	"github.com/beito123/nbt"
)

// DefaultOptions is a default option for leveldb
// You can use at NewWithOptions() and LoadWithOptions()
var DefaultOptions = &opt.Options{
	Filter:      filter.NewBloomFilter(10),
	WriteBuffer: 4 * 1024 * 1024, // 4mb
}

const (
	// LevelDataFile is a location of level.dat
	LevelDataFile = "level.dat"

	// DBPath is a location of level data
	DBPath = "/db"
)

// New returns new LevelDB
// The path is a directory for save
func New(path string) (*LevelDB, error) {
	return NewWithOptions(path, DefaultOptions)
}

// NewWithOptions returns new LevelDB with leveldb options
func NewWithOptions(path string, options *opt.Options) (*LevelDB, error) {
	path = filepath.Clean(path)

	SaveLevelData(path, DefaultProperties)

	db, err := lvldb.OpenFile(filepath.Join(path, DBPath), options)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		Database:   db,
		Format:     &ChunkFormatV100{},
		properties: DefaultProperties,
		chunks:     make(map[uint64]*Chunk),
		mutex:      new(sync.RWMutex),
	}, nil
}

// Load loads a leveldb level
func Load(path string) (*LevelDB, error) {
	return LoadWithOptions(path, DefaultOptions)
}

// LoadWithOptions loads a leveldb level with leveldb options
func LoadWithOptions(path string, options *opt.Options) (*LevelDB, error) {
	path = filepath.Clean(path)

	db, err := lvldb.OpenFile(filepath.Join(path, DBPath), options)
	if err != nil {
		return nil, err
	}

	properties, err := LoadLevelData(path)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		Database:   db,
		Format:     &ChunkFormatV100{},
		properties: properties,
		chunks:     make(map[uint64]*Chunk),
		mutex:      new(sync.RWMutex),
	}, nil
}

// LevelDB is a level format for mcbe
type LevelDB struct {
	Database *lvldb.DB

	Format ChunkFormat

	properties *Properties

	dimension level.Dimension
	chunks    map[uint64]*Chunk

	mutex *sync.RWMutex
}

func (LevelDB) at(x, y int) uint64 {
	return (uint64(uint32(y)) << 32) | uint64(uint32(x))
}

// Name returns name of level
func (lvl *LevelDB) Name() string {
	tag, _ := lvl.Property(TagLevelName) // must

	name, _ := tag.ToString()

	return name
}

// SetName sets the name of level
func (lvl *LevelDB) SetName(name string) {
	lvl.SetProperty(nbt.NewStringTag(TagGameType, name))
}

// GameType returns the default game mode of level
func (lvl *LevelDB) GameType() level.GameType {
	tag, _ := lvl.Property(TagGameType)

	typ, _ := tag.ToInt() // must
	return level.GameType(typ)
}

// SetGameType sets the game mode of level
func (lvl *LevelDB) SetGameType(typ level.GameType) {
	lvl.SetProperty(nbt.NewByteTag(TagGameType, int8(typ)))
}

// Spawn returns the default spawn of level
func (lvl *LevelDB) Spawn() (x, y, z int) {
	spawnX, _ := lvl.Property(TagSpawnX)
	spawnY, _ := lvl.Property(TagSpawnY)
	spawnZ, _ := lvl.Property(TagSpawnZ)

	x, _ = spawnX.ToInt() // must
	y, _ = spawnY.ToInt()
	z, _ = spawnZ.ToInt()

	return x, y, z
}

// SetSpawn sets the default spawn of level
func (lvl *LevelDB) SetSpawn(x, y, z int) {
	lvl.SetProperty(nbt.NewIntTag(TagSpawnX, int32(x)))
	lvl.SetProperty(nbt.NewIntTag(TagSpawnY, int32(x)))
	lvl.SetProperty(nbt.NewIntTag(TagSpawnZ, int32(x)))
}

// Property returns a property of level.dat
func (lvl *LevelDB) Property(name string) (tag nbt.Tag, ok bool) {
	return lvl.properties.Data.Get(name)
}

// SetProperty sets a property
func (lvl *LevelDB) SetProperty(tag nbt.Tag) {
	lvl.properties.Data.Set(tag)
}

// AllProperties returns all properties
func (lvl *LevelDB) AllProperties() *nbt.Compound {
	return lvl.properties.Data
}

// SetAllProperties sets all properties
func (lvl *LevelDB) SetAllProperties(com *nbt.Compound) {
	lvl.properties.Data = com
}

// PropertiesVersion returns properties version
func (lvl *LevelDB) PropertiesVersion() int {
	return lvl.properties.Version
}

// SetPropertiesVersion sets properties version
func (lvl *LevelDB) SetPropertiesVersion(ver int) {
	lvl.properties.Version = ver
}

// Close closes database of leveldb
// You must close after you use the format
// It's should not run other functions after format is closed
func (lvl *LevelDB) Close() error {
	if lvl.Database != nil {
		lvl.mutex.Lock()
		err := lvl.Database.Close()
		lvl.mutex.Unlock()
		return err
	}

	return nil
}

// Dimension return dimension of the level
func (lvl *LevelDB) Dimension() level.Dimension {
	return lvl.dimension
}

// SetDimension set dimension of the level
func (lvl *LevelDB) SetDimension(dimension level.Dimension) {
	lvl.dimension = dimension
}

// LoadChunk loads a chunk.
// If create is enabled, generates a chunk if it doesn't exist
func (lvl *LevelDB) LoadChunk(x, y int, create bool) error {
	if lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: already loaded the chunk")
	}

	exist, err := lvl.HasGeneratedChunk(x, y)
	if err != nil {
		return err
	}

	if !exist && create {
		return lvl.GenerateChunk(x, y)
	}

	chunk, err := lvl.Format.Read(lvl.Database, x, y, lvl.dimension)
	if err != nil {
		return err
	}

	lvl.mutex.Lock()
	lvl.chunks[lvl.at(x, y)] = chunk
	lvl.mutex.Unlock()

	return nil
}

// UnloadChunk unloads a chunk.
func (lvl *LevelDB) UnloadChunk(x, y int) error {
	if !lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: not loaded the chunk")
	}

	lvl.mutex.Lock()
	delete(lvl.chunks, lvl.at(x, y))
	lvl.mutex.Unlock()

	return nil
}

// GenerateChunk generates a chunk and loads
func (lvl *LevelDB) GenerateChunk(x, y int) error {
	if lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: already loaded the chunk")
	}

	chunk := NewChunk(x, y)

	lvl.mutex.Lock()
	lvl.chunks[lvl.at(x, y)] = chunk
	lvl.mutex.Unlock()

	return nil
}

// HasGeneratedChunk returns whether the chunk is generaged
func (lvl *LevelDB) HasGeneratedChunk(x, y int) (bool, error) {
	return lvl.Format.Exist(lvl.Database, x, y, lvl.dimension)
}

// IsLoadedChunk returns weather a chunk is loaded.
func (lvl *LevelDB) IsLoadedChunk(x, y int) bool {
	lvl.mutex.RLock()
	_, ok := lvl.chunks[lvl.at(x, y)]
	lvl.mutex.RUnlock()

	return ok
}

// SaveChunk saves a chunk.
func (lvl *LevelDB) SaveChunk(x, y int) error {
	chunk, ok := lvl.chunk(x, y)
	if !ok {
		return fmt.Errorf("level.leveldb: not loaded the chunk")
	}

	return lvl.Format.Write(lvl.Database, chunk, lvl.dimension)
}

// SaveChunks saves all chunks.
func (lvl *LevelDB) SaveChunks() error {
	for _, chunk := range lvl.chunks {
		err := lvl.SaveChunk(chunk.x, chunk.y)
		if err != nil {
			return err
		}
	}

	return nil
}

func (lvl *LevelDB) chunk(x, y int) (*Chunk, bool) {
	lvl.mutex.RLock()
	chunk, ok := lvl.chunks[lvl.at(x, y)]
	lvl.mutex.RUnlock()

	return chunk, ok
}

// Chunk returns a chunk.
// If a chunk is not loaded, it will be loaded
func (lvl *LevelDB) Chunk(x, y int) (level.Chunk, error) {
	if lvl.IsLoadedChunk(x, y) {
		chunk, ok := lvl.chunk(x, y)
		if !ok {
			return nil, errors.New("couldn't find the chunk")
		}

		return chunk, nil
	}

	err := lvl.LoadChunk(x, y, false)
	if err != nil {
		return nil, err
	}

	chunk, ok := lvl.chunk(x, y)
	if !ok {
		return nil, errors.New("couldn't find the chunk")
	}

	return chunk, nil
}

// LoadedChunks returns loaded chunks.
func (lvl *LevelDB) LoadedChunks() []level.Chunk {
	result := make([]level.Chunk, len(lvl.chunks))

	lvl.mutex.RLock()

	count := 0
	for _, chunk := range lvl.chunks {
		result[count] = chunk
		count++
	}

	lvl.mutex.RUnlock()

	return result
}

// NewRawBlockState returns new RawBlockState
func NewRawBlockState(name string, value int) *RawBlockState {
	return &RawBlockState{
		name:  strings.ToLower(name),
		value: value,
	}
}

// FromRawBlockState returns new RawBlockState
func FromRawBlockState(bs level.BlockState) (*RawBlockState, error) {
	name, meta, ok := bs.ToBlockNameMeta()
	if !ok {
		return nil, fmt.Errorf("level.leveldb: usable to convert from %s to RawBlockState", bs.Name())
	}

	return NewRawBlockState(name, meta), nil
}

// RawBlockState is a raw block information
type RawBlockState struct {
	name  string
	value int
}

// Name returns block name
func (block *RawBlockState) Name() string {
	return block.name
}

// Value returns block value
func (block *RawBlockState) Value() int {
	return block.value
}

// Equal returns whether block is equal b
func (block *RawBlockState) Equal(b *RawBlockState) bool {
	return block.name == b.name && block.value == b.value
}

// ToBlockNameProperties returns block name and properties
// If it's not supported, returns false for ok
func (block *RawBlockState) ToBlockNameProperties() (name string, properties map[string]string, ok bool) {
	return block.name, make(map[string]string), true
}

// ToBlockNameMeta returns block name and meta
// If it's not supported, returns false for ok
func (block *RawBlockState) ToBlockNameMeta() (name string, meta int, ok bool) {
	return block.name, block.value, true
}

// ToBlockIDMeta returns block id and meta
// If it's not supported, returns false for ok
func (block *RawBlockState) ToBlockIDMeta() (id int, meta int, ok bool) {
	return 0, 0, false
}
