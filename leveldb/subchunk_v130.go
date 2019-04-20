package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"
	"math"

	"github.com/beito123/level/binary"
	"github.com/beito123/level/block"
	"github.com/beito123/level/util"
	"github.com/beito123/nbt"
)

type SubChunkFormatV130 struct {
	RuntimeIDList map[int]block.BlockData
}

func (format *SubChunkFormatV130) Version() byte {
	return 0x08
}

func (format *SubChunkFormatV130) Read(y byte, b []byte) (*SubChunk, error) {
	sub := NewSubChunk(y)

	stream := binary.NewStreamBytes(b)

	ver, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	if ver != format.Version() {
		return nil, fmt.Errorf("level.leveldb: unsupported version: version %d", ver)
	}

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

	return nil, nil
}

func (format *SubChunkFormatV130) ReadBlockStorage(stream *binary.Stream) (*BlockStorage, error) {
	storage := NewBlockStorage()

	flags, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	bitsPerBlock := (flags << 1)
	isRuntime := (flags & 0x01) != 0

	if bitsPerBlock > 16 {
		return nil, fmt.Errorf("level.leveldb: unsupported bits per block, wants 1-16 bits")
	}

	mask := uint16((1 << uint(bitsPerBlock)) - 1)

	wordBits := 8 * 4 // 1byte * 4
	blocksPerWord := float64(wordBits) / float64(bitsPerBlock)

	wordCount := util.CeilInt(float64(BlockStorageSize) / math.Ceil(blocksPerWord))

	count := 0
	for i := 0; i < wordCount; i++ {
		word, err := stream.Int()
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(blocksPerWord); j++ {
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
	} else { // nbt format
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
	}

	return storage, nil
}
