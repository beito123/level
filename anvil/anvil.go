package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	"github.com/beito123/binary"
	"github.com/beito123/level"
	"github.com/beito123/nbt"
)

func NewAnvil(path string) (*Anvil, error) {
	loader, err := NewRegionLoader(path, RegionFileAnvil)
	if err != nil {
		return nil, err
	}

	return &Anvil{
		loader:  loader,
		regions: make(map[uint64]*Region),
	}, nil
}

// Anvil is a level format
// It often is used for minecraft java edition and server world
type Anvil struct {
	loader  *RegionLoader
	regions map[uint64]*Region
}

// getRegion returns a region with xy
// If the Region doesn't exist, ok is false
func (lvl *Anvil) getRegion(x, y int) (r *Region, ok bool) {
	r, ok = lvl.regions[lvl.toIndex(x, y)]

	return r, ok
}

// toCC returns id for container from region coordinate
func (Anvil) toIndex(x, y int) uint64 {
	return uint64(int64(x)<<32 | int64(y))
}

// chunkToRegion returns region coordinate from chunk coordinate
func (Anvil) chunkToRegion(x, y int) (cx, cy int) {
	return x >> 5, y >> 5
}

// LoadChunk loads a chunk.
// If create is true, generates a chunk.
func (lvl *Anvil) LoadChunk(x, y int, create bool) bool {
	return false
}

// UnloadChunk unloads a chunk.
// If safe is false, unloads even when players are there.
func (lvl *Anvil) UnloadChunk(x, y int, safe bool) bool {
	return false
}

// IsLoadedChunk returns weather a chunk is loaded.
func (lvl *Anvil) IsLoadedChunk(x, y int) bool {
	return false
}

// SaveChunk saves a chunk.
func (lvl *Anvil) SaveChunk(x, y int) bool {
	return false
}

// SaveChunks saves all chunks.
func (lvl *Anvil) SaveChunks() {

}

// Chunk returns a chunk.
// If it's not loaded, loads the chunks.
// If create is true, generates a chunk.
func (lvl *Anvil) Chunk(x, y int, create bool) level.Chunk {
	return nil
}

// Chunks retuns loaded chunks.
func (lvl *Anvil) Chunks() []level.Chunk {
	return nil
}

// NewChunk returns new Chunk
func NewChunk(x, y int, subFormat SubChunkFormat) (*Chunk, error) {
	return &Chunk{
		x:              x,
		y:              y,
		SubChunkFormat: subFormat,
	}, nil
}

// ReadChunk returns new Chunk with CompoundTag
func ReadChunk(x, y int, b []byte) (*Chunk, error) {
	stream := nbt.NewStreamBytes(nbt.BigEndian, b)

	tag, err := stream.ReadTag()
	if err != nil {
		return nil, err
	}

	com, ok := tag.(*nbt.Compound)
	if !ok {
		return nil, fmt.Errorf("level.anvil.region: expected to be CompoundTag, but it passed different tag(%sTag)", nbt.GetTagName(tag.ID()))
	}

	subFormat := &SubChunkFormatB{}
	if com.Has("DataVersion") { // introduced v1.13
		ver, err := com.GetInt("DataVersion")
		if err != nil {
			return nil, err
		}

		switch ver {
		default:
			subFormat = &SubChunkFormatB{}
		}
	}

	chunk, err := NewChunk(x, y, subFormat)
	if err != nil {
		return nil, err
	}

	fmt.Printf("x: %d, y: %d \n", x, y)

	err = chunk.Load(com)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// Chunk is
type Chunk struct {
	x int
	y int

	lastUpdate    int64
	inhabitedTime int64
	biomes        []int
	subChunks     []*SubChunk

	SubChunkFormat SubChunkFormat
}

// X returns x coordinate
func (chunk *Chunk) X() int {
	return chunk.x
}

// Y returns y coordinate
func (chunk *Chunk) Y() int {
	return chunk.y
}

func (chunk *Chunk) SubChunks() []*SubChunk {
	return chunk.subChunks
}

func (chunk *Chunk) Load(tag *nbt.Compound) error {

	com, err := tag.GetCompound("Level")
	if err != nil {
		return err
	}

	// Biomes
	biomes, err := com.GetIntArray("Biomes")
	if err != nil {
		return err
	}

	chunk.biomes = make([]int, len(biomes))
	for i, biome := range biomes {
		chunk.biomes[i] = int(biome)
	}

	// Subchunks
	sections, err := com.GetList("Sections")
	if err != nil {
		return err
	}

	chunk.subChunks = make([]*SubChunk, 16)
	for _, entry := range sections {
		sec, ok := entry.(*nbt.Compound)
		if !ok {
			return fmt.Errorf("couldn't convert to *Compound")
		}

		sub, err := chunk.SubChunkFormat.Read(sec)
		if err != nil {
			return err
		}

		chunk.subChunks[sub.Y] = sub
	}

	return nil
}

type SubChunkFormat interface {
	Read(tag *nbt.Compound) (*SubChunk, error)
}

type SubChunkFormatA struct { // before v1.13
}

func (SubChunkFormatA) Read(tag *nbt.Compound) (*SubChunk, error) {
	return nil, nil
}

type SubChunkFormatB struct { // after v1.13
}

func (SubChunkFormatB) Read(tag *nbt.Compound) (*SubChunk, error) {
	y, err := tag.GetByte("Y")
	if err != nil {
		return nil, err
	}

	sub := &SubChunk{
		Y: y,
	}

	blockData, err := tag.GetLongArray("BlockStates")
	if err != nil {
		return nil, err
	}

	blockCount := 16 * 16 * 16 // 4096

	longLen := binary.LongSize * 8 // 1bytes = 8bits // 64

	bit := (len(blockData) * longLen) / blockCount      // bits per block
	perBlock := (len(blockData) * longLen) / blockCount // bits per block
	mask := uint16((1 << uint(bit)) - 1)                // returns 4bits -> 0b1111

	fmt.Printf("bits: %02b\n", mask)

	sub.Blocks = make([]uint16, blockCount)

	var count int
	/*for _, data := range blockData {
		for i := 0; i < (longLen / bit); i++ {
			id := uint16(data>>uint(bit*i)) & uint16(mask)

			sub.Blocks[count] = id

			count++
		}
	}*/

	var c2 = 100
	for i := 0; i < blockCount; i++ {
		index := (perBlock * i) / longLen
		off := (perBlock * i) % longLen

		data := uint16(uint64(blockData[index])>>uint(off)) & mask

		left := ((off + perBlock) - longLen)
		if left > 0 {
			m := uint64((1 << uint(left)) - 1)
			data = ((uint16(uint64(blockData[index+1])&m) << uint(perBlock-left)) | data) & mask

			if c2 <= 10 {
				fmt.Printf("base: %064b, \nsub:  %064b\n", uint64(blockData[index]), uint64(blockData[index+1]))
				fmt.Printf("index: %d, offset: %d, left: %d \n", index, off, left)
				fmt.Printf("from: %05b, to: %05b\n", uint16(uint64(blockData[index])>>uint(off))&mask, data)

				c2++
			}

		}

		sub.Blocks[count] = data

		count++
	}

	palettes, err := tag.GetList("Palette")
	if err != nil {
		return nil, err
	}

	sub.Palette = make([]*BlockState, len(palettes))

	for i, entry := range palettes {
		pac, ok := entry.(*nbt.Compound)
		if !ok {
			return nil, fmt.Errorf("couldn't convert to *Compound")
		}

		bName, err := pac.GetString("Name")
		if err != nil {
			return nil, err
		}

		sub.Palette[i] = &BlockState{
			Name: bName,
		}
	}

	return sub, nil
}

func (SubChunkFormatB) Write(tag *nbt.Compound) (*SubChunk, error) {
	return nil, nil
}

// SubChunk is what is divided a chunk horizontally in sixteen
type SubChunk struct {
	Y byte

	Palette    []*BlockState
	Blocks     []uint16
	BlockLight []byte
	SkyLight   []byte
}

type BlockState struct {
	Name       string
	Properties map[string]string
}
