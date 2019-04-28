package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beito123/nbt"

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

func (f Finalization) ID() byte {
	switch f {
	case NotGenerated:
		return 0
	case NotSpawnMobs:
		return 1
	case Generated, Unsupported:
		return 2
	}

	return 0
}

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
	x             int
	y             int
	subChunks     []*SubChunk
	heightMap     []uint16
	biomes        []byte
	entities      []*nbt.Compound
	blockEntities []*nbt.Compound

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

func (chunk *Chunk) atData2D(x, y int) int {
	return y*16 + x
}

// Height returns the height of the highest block at chunk coordinate
func (chunk *Chunk) Height(x, y int) uint16 {
	return chunk.heightMap[chunk.atData2D(x, y)]
}

// Biome returns biome
func (chunk *Chunk) Biome(x, y int) byte {
	return chunk.biomes[chunk.atData2D(x, y)]
}

// SetBiome set biome
func (chunk *Chunk) SetBiome(x, y int, biome byte) {
	chunk.biomes[chunk.atData2D(x, y)] = biome
}

// Entities returns entities of nbt data
func (chunk *Chunk) Entities() []*nbt.Compound {
	return chunk.entities
}

// SetEntities set entities of nbt data
func (chunk *Chunk) SetEntities(entities []*nbt.Compound) {
	chunk.entities = entities
}

// BlockEntities returns block entities of nbt data
func (chunk *Chunk) BlockEntities() []*nbt.Compound {
	return chunk.blockEntities
}

// SetBlockEntities set block entities of nbt data
func (chunk *Chunk) SetBlockEntities(entities []*nbt.Compound) {
	chunk.blockEntities = entities
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
		sub = NewSubChunk(byte(y / 16))
	}

	err := sub.SetBlock(x, y&15, z, index, bs)

	if err != nil {
		return err
	}

	if chunk.Finalization == NotGenerated {
		chunk.Finalization = NotSpawnMobs
	}

	return nil
}

const (
	TagData2D         = 45
	TagData2DLegacy   = 46
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
	Read(db *lvldb.DB, x, y int, dimension level.Dimension) (*Chunk, error)
	Write(db *lvldb.DB, chunk *Chunk, dimension level.Dimension) error
	Exist(db *lvldb.DB, x, y int, dimension level.Dimension) (bool, error)
}

const (
	SubChunkVersionV120  = 0
	SubChunkVersionV1213 = 1
	SubChunkVersionV130  = 8
)

// ChunkFormatV100 is a chunk format v1.0.0 or after
type ChunkFormatV100 struct {
	// SubChunkVersion is used a format when it writes a chunk
	SubChunkVersion int

	DisabledData2D      bool
	DisabledEntity      bool
	DisabledBlockEntity bool
}

// Read reads a chunk
func (format *ChunkFormatV100) Read(db *lvldb.DB, x, y int, dimension level.Dimension) (*Chunk, error) {
	chunk := NewChunk(x, y)
	chunk.DefaultBlock = NewRawBlockState("minecraft:air", 0)

	// Exist

	exist, err := format.Exist(db, x, y, dimension)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("level.leveldb: the chunk isn't generated")
	}

	// Finalization

	stateKey := format.getChunkKey(x, y, dimension, TagFinalizedState, -1)

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

	// Read Data2D
	if !format.DisabledData2D {
		data2dKey := format.getChunkKey(x, y, dimension, TagData2D, -1)

		hasData2D, err := db.Has(data2dKey, nil)
		if err != nil {
			return nil, err
		}

		if hasData2D { // sometimes a chunk hasn't entities
			b, err := db.Get(data2dKey, nil)
			if err != nil {
				return nil, err
			}

			ln := len(b)
			heightMapLen := 512
			biomesLen := 256

			if ln < heightMapLen {
				return nil, fmt.Errorf("level.leveldb: not enough bytes for HeightMap")
			}

			rawHeightMap := b[:512]

			chunk.heightMap = make([]uint16, 16*16)
			for i := 0; i < len(chunk.heightMap); i++ {
				chunk.heightMap[i] = binary.ReadLUShort(rawHeightMap[i*2 : i*2+2])
			}

			if ln < biomesLen {
				return nil, fmt.Errorf("level.leveldb: not enough bytes for Biomes")
			}

			chunk.biomes = b[heightMapLen : heightMapLen+biomesLen]
		}
	}

	// Read Entity
	if !format.DisabledEntity {
		entityKey := format.getChunkKey(x, y, dimension, TagEntity, -1)

		hasEntity, err := db.Has(entityKey, nil)
		if err != nil {
			return nil, err
		}

		if hasEntity { // sometimes a chunk hasn't entities
			b, err := db.Get(entityKey, nil)
			if err != nil {
				return nil, err
			}

			chunk.entities, err = format.ReadCompounds(b)
			if err != nil {
				return nil, err
			}
		}
	}

	// Read BlockEntity
	if !format.DisabledBlockEntity {
		blockEntityKey := format.getChunkKey(x, y, dimension, TagBlockEntity, -1)

		hasBlockEntity, err := db.Has(blockEntityKey, nil)
		if err != nil {
			return nil, err
		}

		if hasBlockEntity { // sometimes a chunk hasn't block entities
			b, err := db.Get(blockEntityKey, nil)
			if err != nil {
				return nil, err
			}

			chunk.blockEntities, err = format.ReadCompounds(b)
			if err != nil {
				return nil, err
			}
		}
	}

	return chunk, nil
}

// Write writes a chunk
func (format *ChunkFormatV100) Write(db *lvldb.DB, chunk *Chunk, dimension level.Dimension) error {
	if chunk.Finalization != Unsupported {
		stateKey := format.getChunkKey(chunk.X(), chunk.X(), dimension, TagFinalizedState, -1)

		err := db.Put(stateKey, []byte{chunk.Finalization.ID()}, nil)
		if err != nil {
			return err
		}
	}

	// Write subchunks
	for _, sub := range chunk.SubChunks() {
		b, err := format.WriteSubChunk(sub)
		if err != nil {
			return err
		}

		key := format.getChunkKey(chunk.x, chunk.y, dimension, TagSubChunkPrefix, int(sub.Y))

		err = db.Put(key, b, nil)
		if err != nil {
			return err
		}
	}

	if !format.DisabledEntity {
		b, err := format.WriteCompounds(chunk.entities)
		if err != nil {
			return err
		}

		err = db.Put(format.getChunkKey(chunk.x, chunk.y, dimension, TagEntity, -1), b, nil)
		if err != nil {
			return err
		}
	}

	if !format.DisabledBlockEntity {
		b, err := format.WriteCompounds(chunk.blockEntities)
		if err != nil {
			return err
		}

		err = db.Put(format.getChunkKey(chunk.x, chunk.y, dimension, TagBlockEntity, -1), b, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Exist returns whether a chunk is generated
func (format *ChunkFormatV100) Exist(db *lvldb.DB, x, y int, dimension level.Dimension) (bool, error) {
	return db.Has(format.getChunkKey(x, y, dimension, TagVersion, -1), nil)
}

// ReadSubChunk reads a subchunk from bytes b
func (format *ChunkFormatV100) ReadSubChunk(y byte, b []byte) (sub *SubChunk, err error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("level.leveldb: not enough bytes")
	}

	ver := b[0]

	switch ver {
	case SubChunkVersionV120, 2, 3, 4, 5, 6, 7: // v1.2 or before
		// TODO: support old format
		return nil, fmt.Errorf("level.leveldb: unsupported old subchunk format")
	case SubChunkVersionV1213, SubChunkVersionV130: // Palettized format // 1.2.13 or after
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

// WriteSubChunk reads a subchunk from bytes b
func (format *ChunkFormatV100) WriteSubChunk(sub *SubChunk) (b []byte, err error) {
	if format.SubChunkVersion == SubChunkVersionV120 { // TODO: support
		return nil, fmt.Errorf("level.leveldb: unsupported old subchunk format")
	}

	switch format.SubChunkVersion {
	case SubChunkVersionV120, 2, 3, 4, 5, 6, 7:
		return nil, fmt.Errorf("level.leveldb: unsupported old subchunk format")
	case SubChunkVersionV1213, SubChunkVersionV130:
		subFormat := &SubChunkFormatV1213{
			OldFormat: format.SubChunkVersion == SubChunkVersionV1213,
		}

		b, err = subFormat.Write(sub)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("level.leveldb: unsupported subchunk version %d", format.SubChunkVersion)
	}

	return b, nil
}

// ReadCompounds reads compounds
func (format *ChunkFormatV100) ReadCompounds(b []byte) ([]*nbt.Compound, error) {
	var list []*nbt.Compound

	stream := nbt.NewStreamBytes(nbt.LittleEndian, b)
	for i := 0; i < 65536; i++ { // Limit for infinite loop
		tag, err := stream.ReadTag()
		if err != nil {
			ioutil.WriteFile("./test.nbt", b, os.ModePerm)
			return nil, err
		}

		com, ok := tag.(*nbt.Compound)
		if !ok {
			return nil, fmt.Errorf("level.leveldb: couldn't convert to nbt.Compound")
		}

		list = append(list, com)

		if stream.Stream.Len() == 0 {
			break
		}
	}

	return list, nil
}

// WriteCompounds writes compounds
func (format *ChunkFormatV100) WriteCompounds(tags []*nbt.Compound) ([]byte, error) {
	stream := nbt.NewStream(nbt.LittleEndian)
	for _, com := range tags {
		err := stream.WriteTag(com)
		if err != nil {
			return nil, err
		}
	}

	return stream.Bytes(), nil
}

func (format *ChunkFormatV100) toDimensionID(dimension level.Dimension) int {
	switch dimension {
	case level.OverWorld:
		return 0
	case level.Nether:
		return 10
	case level.TheEnd:
		return 20
	}

	return 0
}

func (format *ChunkFormatV100) fromDimensionID(id int) level.Dimension {
	switch id {
	case 0:
		return level.OverWorld
	case 10:
		return level.Nether
	case 20:
		return level.TheEnd
	}

	return level.Unknown
}

func (format *ChunkFormatV100) getChunkKey(x int, y int, dimension level.Dimension, tag byte, sid int) []byte {
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

	dimID := format.toDimensionID(dimension)

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
			byte(dimID),
			byte(dimID >> 8),
			byte(dimID >> 16),
			byte(dimID >> 24),
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
			byte(dimID),
			byte(dimID >> 8),
			byte(dimID >> 16),
			byte(dimID >> 24),
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
