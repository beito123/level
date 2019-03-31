package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	"github.com/beito123/level/block"

	"github.com/beito123/binary"
	"github.com/beito123/nbt"
)

// SubChunkFormat is a subchunk format for a version
type SubChunkFormat interface {
	Read(tag *nbt.Compound) (*SubChunk, error)
}

// SubChunkFormatV112 is a subchunk format for v1.12 and before
type SubChunkFormatV112 struct {
}

func (SubChunkFormatV112) Read(tag *nbt.Compound) (*SubChunk, error) {
	y, err := tag.GetByte("Y")
	if err != nil {
		return nil, err
	}

	sub := &SubChunk{
		Y: y,
	}

	// Blocks

	sub.Palette = []*BlockState{
		0: &BlockState{
			Name: "minecraft:air",
		},
	}

	blocks, err := tag.GetByteArray("Blocks")
	if err != nil {
		return nil, err
	}

	data, err := tag.GetByteArray("Data")
	if err != nil {
		return nil, err
	}

	blockCount := 16 * 16 * 16 // 4096

	sub.Blocks = make([]uint16, blockCount)

	for i := range sub.Blocks {
		state := &BlockState{
			IsOld:   true,
			OldID:   blocks[i],
			OldMeta: ToNibble(data, i),
		}

		index := -1
		for ind, val := range sub.Palette { // find palette
			if val.Equal(state) {
				index = ind
				break
			}
		}

		if index == -1 {
			index = len(sub.Palette) // next index
			sub.Palette = append(sub.Palette, state)
		}

		sub.Blocks[i] = uint16(index)
	}

	// BlockLight

	sub.BlockLight, err = tag.GetByteArray("BlockLight")
	if err != nil {
		return nil, err
	}

	// SkyLight

	sub.SkyLight, err = tag.GetByteArray("SkyLight")
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// SubChunkFormatV113 is a subchunk format for v1.13 and after
type SubChunkFormatV113 struct { // after v1.13
}

func (SubChunkFormatV113) Read(tag *nbt.Compound) (*SubChunk, error) {
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

	//fmt.Printf("bits: %02b\n", mask)

	sub.Blocks = make([]uint16, blockCount)

	var count int

	for i := 0; i < blockCount; i++ {
		index := (perBlock * i) / longLen
		off := (perBlock * i) % longLen

		data := uint16(uint64(blockData[index])>>uint(off)) & mask

		left := ((off + perBlock) - longLen)
		if left > 0 {
			m := uint64((1 << uint(left)) - 1)
			data = ((uint16(uint64(blockData[index+1])&m) << uint(perBlock-left)) | data) & mask
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

	// BlockLight

	sub.BlockLight, err = tag.GetByteArray("BlockLight")
	if err != nil {
		return nil, err
	}

	// SkyLight

	sub.SkyLight, err = tag.GetByteArray("SkyLight")
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (SubChunkFormatV113) Write(sub *SubChunk) (*nbt.Compound, error) {
	return nil, nil
}

// NewSubChunk returns new subchunk
func NewSubChunk(y byte) *SubChunk {
	return &SubChunk{
		Y:          y,
		Palette:    make([]*BlockState, 0),
		Blocks:     make([]uint16, 2048),
		BlockLight: make([]byte, 2048),
		SkyLight:   make([]byte, 2048),
	}
}

// SubChunk is what is divided a chunk horizontally in sixteen
type SubChunk struct {
	Y byte

	Palette    []*BlockState
	Blocks     []uint16
	BlockLight []byte
	SkyLight   []byte
}

// At returns index from subchunk coordinates
// xyz need to be more 0 and less 15
func (SubChunk) At(x, y, z int) int {
	return y<<8 | z<<4 | x
}

// Vaild vailds subchunk coordinates
func (SubChunk) Vaild(x, y, z int) error {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return fmt.Errorf("invail coordinate")
	}

	return nil
}

// AtBlock returns block name at the subchunk coordinates
func (sub *SubChunk) AtBlock(x, y, z int) (*BlockState, error) {
	id, err := sub.BlockIndex(x, y, z)
	if err != nil {
		return nil, err
	}

	if int(id) >= len(sub.Palette) {
		return nil, fmt.Errorf("couldn't find a palette for the block")
	}

	return sub.Palette[id], nil
}

// BlockIndex returns a index for block id
func (sub *SubChunk) BlockIndex(x, y, z int) (uint16, error) {
	err := sub.Vaild(x, y, z)
	if err != nil {
		return 0, err
	}

	return sub.Blocks[sub.At(x, y, z)], nil
}

// AtBlockLight returns a blocklight at the subchunk coordinates
func (sub *SubChunk) AtBlockLight(x, y, z int) (byte, error) {
	err := sub.Vaild(x, y, z)
	if err != nil {
		return 0, err
	}

	return ToNibble(sub.BlockLight, sub.At(x, y, z)), nil
}

// AtSkyLight returns a skylight at the subchunk coordinates
func (sub *SubChunk) AtSkyLight(x, y, z int) (byte, error) {
	err := sub.Vaild(x, y, z)
	if err != nil {
		return 0, err
	}

	return ToNibble(sub.SkyLight, sub.At(x, y, z)), nil
}

/*
type BlockRawData struct {
	Name       string
	Properties map[string]string
}*/

type BlockState struct {
	Name       string
	Properties map[string]string

	IsOld   bool
	OldID   byte
	OldMeta byte
}

func (bs *BlockState) ToBlockData() *block.BlockData {
	if bs.IsOld {
		return block.FromBlockID(int(bs.OldID), int(bs.OldMeta))
	}

	return &block.BlockData{
		Name:       bs.Name,
		Properties: bs.Properties,
	}
}

func (bs *BlockState) Equal(sub *BlockState) bool {
	if bs.IsOld {
		return bs.OldID == sub.OldID && bs.OldMeta == sub.OldMeta
	}

	if bs.Name != sub.Name {
		return false
	}

	if len(bs.Properties) != len(sub.Properties) {
		return false
	}

	for k, v := range bs.Properties {
		val, ok := sub.Properties[k]
		if !ok {
			return false
		} else if v != val {
			return false
		}
	}

	return true
}

// ToNibble returns a nibble data from []byte by index
func ToNibble(b []byte, index int) byte {
	data := b[index/2]

	if (index % 2) != 0 { // 0b111100000
		return (data >> 4) & 0x0F
	}

	return data & 0x0F // 0b00001111
}
