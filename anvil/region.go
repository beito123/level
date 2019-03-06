package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beito123/nbt"

	"github.com/beito123/level/binary"
	"github.com/beito123/level/util"
)

// NewRegionLoader returns new RegionLoader
func NewRegionLoader(path string) (*RegionLoader, error) {
	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !f.IsDir() { // = a file path
		return nil, errors.New("it couldn't use a file path")
	}

	return &RegionLoader{
		path: path,
	}, nil
}

// RegionLoader controls a region file on dir and pass Region to load a region
type RegionLoader struct {
	path string
}

func (RegionLoader) toRegionFile(x, y int) string { // Should I write to be able to change?
	return "r." + strconv.Itoa(x) + "." + strconv.Itoa(y) + ".mca"
}

func (rl *RegionLoader) createRegion(x, y int) (*Region, error) {
	return nil, nil
}

// LoadRegion loads a region
func (rl *RegionLoader) LoadRegion(x, y int, create bool) (*Region, error) {
	path := rl.path + "/" + rl.toRegionFile(x, y)

	if !util.ExistFile(path) {
		if !create {
			return nil, fmt.Errorf("level.anvil: couldn't find the region (x: %d, y: %d)", x, y)
		}

		return rl.createRegion(x, y)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reg := &Region{
		X: x,
		Y: y,
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = reg.Load(b)
	if err != nil {
		return nil, err
	}

	return reg, nil
}

const (
	// ChunkCount is the number of chunks in a region file
	ChunkCount = 32 * 32 // 1024

	// Sector is a data sector for region format
	Sector = 4 * ChunkCount

	// LocationsBytes is a offset for chunks location in a region file
	LocationsBytes = Sector // 4 = (offset 3bytes, count 1byte)

	// TimestampsBytes is a offset for chunks's timestamp in a region file
	TimestampsBytes = Sector // 4 = timestamp 4 bytes

	// InformationSector is a information sector for chunk data
	InformationSector = LocationsBytes + TimestampsBytes
)

// Region is a section had 32x32 chunks
type Region struct {
	X int
	Y int

	Chunks map[uint64]*Chunk
}

// Load loads a region from region file bytes
func (reg *Region) Load(b []byte) error {
	if len(b) < InformationSector {
		return fmt.Errorf("level.anvil.region: the region bytes isn't not enough")
	}

	reg.Chunks = make(map[uint64]*Chunk, ChunkCount)

	stream := binary.NewStreamBytes(b[:LocationsBytes+TimestampsBytes])

	// Read Locations
	locations := make([]*Location, ChunkCount)

	for i := 0; i < len(locations); i++ {
		offset, err := stream.Triad() // location of chunk data
		if err != nil {
			return err
		}

		count, err := stream.Byte() // lenght of chunk data
		if err != nil {
			return err
		}

		locations[i] = &Location{
			Off:   offset,
			Count: count,
		}
	}

	// Read timestamps
	timestamps := make([]int32, ChunkCount)

	for i := 0; i < len(timestamps); i++ {
		stamp, err := stream.Int()
		if err != nil {
			return err
		}

		timestamps[i] = stamp
	}

	stream.Reset() // the bytes won't be use

	// Read chunk data
	for _, locat := range locations {
		off := int(locat.Off) * Sector
		ln := int(locat.Count) * Sector

		if off == 0 && ln == 0 { // It haven't generated yet
			continue
		} else if off < 2 {
			return fmt.Errorf("level.anvil.region: invaild offset")
		}

		stream = binary.NewStreamBytes(b[off : off+ln]) // chunk data + pads

		realLen, err := stream.Int() // ln = readLen + pad(realLen % 4096bytes)
		if err != nil {
			return err
		}

		stream.Skip(1) // compression type, but we won't use

		nstream, err := nbt.FromBytes(stream.Get(int(realLen)), nbt.BigEndian) // NBT Data
		if err != nil {
			return err
		}

		tag, err := nstream.ReadTag()
		if err != nil {
			return err
		}

		com, ok := tag.(*nbt.Compound)
		if !ok {
			return fmt.Errorf("level.anvil.region: expected to be CompoundTag, but it passed different tag(%sTag)", nbt.GetTagName(tag.ID()))
		}

		chunk, err := ReadChunk(com)
		if err != nil {
			return err
		}

		reg.Chunks[reg.toKey(chunk.X(), chunk.Y())] = chunk
	}

	return nil
}

func (Region) toKey(x, y int) uint64 {
	return uint64(int64(x)<<32 | int64(y))
}

// Save saves region data, returns bytes for a region file
func (reg *Region) Save() ([]byte, error) {
	// TODO: write

	return nil, nil
}

// Location is a location info for chunk data
type Location struct {
	Off   binary.Triad
	Count byte
}
