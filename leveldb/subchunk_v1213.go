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
	"github.com/beito123/level/util"
	"github.com/beito123/nbt"
)

// SubChunkFormatV1213 is a subchunk formatter v1.2.13 or after
type SubChunkFormatV1213 struct {
	RuntimeIDList map[int]*BlockState
}

func (format *SubChunkFormatV1213) Read(y byte, b []byte) (*SubChunk, error) {
	sub := NewSubChunk(y)

	stream := binary.NewStreamBytes(b)

	ver, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	switch ver {
	case 1: // v1.2.13
		storage, err := format.ReadBlockStorage(stream)
		if err != nil {
			return nil, err
		}

		sub.Storages = append(sub.Storages, storage)
	case 8:
		numStorage, err := stream.Byte()
		if err != nil {
			return nil, err
		}

		for i := 0; i < int(numStorage); i++ {
			storage, err := format.ReadBlockStorage(stream)
			if err != nil {
				return nil, err
			}

			sub.Storages = append(sub.Storages, storage)
		}
	default:
		return nil, fmt.Errorf("level.leveldb: unsupported version: version %d", ver)
	}

	return sub, nil
}

// ReadBlockStorage reads a block storage
func (format *SubChunkFormatV1213) ReadBlockStorage(stream *binary.Stream) (*BlockStorage, error) {
	storage := NewBlockStorage()

	flags, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	bitsPerBlock := flags >> 1
	isRuntime := (flags & 0x01) != 0

	if bitsPerBlock > 16 {
		return nil, fmt.Errorf("level.leveldb: unsupported bits per block, wants 1-16 bits")
	}

	mask := uint16((1 << uint(bitsPerBlock)) - 1)

	wordBits := 8 * 4 // 1byte * 4
	blocksPerWord := wordBits / int(bitsPerBlock)

	wordCount := util.CeilInt(float64(BlockStorageSize) / float64(blocksPerWord))

	count := 0
	for i := 0; i < wordCount; i++ {
		word, err := stream.LInt()
		if err != nil {
			return nil, err
		}

		for j := 0; j < blocksPerWord && count < BlockStorageSize; j++ {
			id := uint16(word>>uint(j*int(bitsPerBlock))) & mask

			storage.Blocks[count] = id

			count++
		}
	}

	paletteSize, err := stream.LInt()
	if err != nil {
		return nil, err
	}

	if isRuntime { // I don't have binary data, please give me if you have
		return nil, fmt.Errorf("level.leveldb: unsupported runtime id")
	}

	nbtStream := nbt.NewStreamBytes(nbt.LittleEndian, stream.Bytes())

	for i := 0; i < int(paletteSize); i++ {
		tag, err := nbtStream.ReadTag()
		if err != nil {
			return nil, err
		}

		com, ok := tag.(*nbt.Compound)
		if !ok {
			return nil, fmt.Errorf("level.leveldb: unexpected tag")
		}

		name, err := com.GetString("name")
		if err != nil {
			return nil, err
		}

		val, err := com.GetInt("val")
		if err != nil {
			return nil, err
		}

		state := NewBlockState(name, int(val))

		storage.Palettes = append(storage.Palettes, state)
	}

	stream.Skip(nbtStream.Stream.Off())

	return storage, nil
}
