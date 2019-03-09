package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beito123/level/binary"
	"github.com/beito123/level/util"
)

var (
	// RegionFileMCRegion returns a region file name for mcregion
	RegionFileMCRegion = func(x, y int) string {
		return "r." + strconv.Itoa(x) + "." + strconv.Itoa(y) + ".mcr"
	}

	// RegionFileAnvil returns a region file name for mcregion
	RegionFileAnvil = func(x, y int) string {
		return "r." + strconv.Itoa(x) + "." + strconv.Itoa(y) + ".mca"
	}
)

//

// NewRegionLoader returns new RegionLoader
// You can set RegionFileMCRegion and RegionFileAnvil to tofile
func NewRegionLoader(path string, tofile func(x, y int) string) (*RegionLoader, error) {
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
		path:         path,
		ToRegionFile: tofile,
	}, nil
}

// RegionLoader controls a region file on dir and pass Region to load a region
type RegionLoader struct {
	path string

	ToRegionFile func(x, y int) string
}

// LoadRegion loads a region
func (rl *RegionLoader) LoadRegion(x, y int, create bool) (*Region, error) {
	path := util.To(rl.path, rl.ToRegionFile(x, y))

	if !util.ExistFile(path) {
		if !create {
			return nil, fmt.Errorf("level.anvil: couldn't find the region (x: %d, y: %d)", x, y)
		}

		return NewRegion(x, y), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reg := NewRegion(x, y)

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

// SaveRegion saves a region as a file
func (rl *RegionLoader) SaveRegion(reg *Region) error {
	b, err := reg.Save()
	if err != nil {
		return err
	}

	path := util.To(rl.path, rl.ToRegionFile(reg.X, reg.Y))

	return ioutil.WriteFile(path, b, os.ModePerm)
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

const (
	// CompressionGZip compresses chunk data with gzip for chunk data
	CompressionGZip = iota + 1

	// CompressionZlib compresses chunk data with zlib for chunk data
	CompressionZlib
)

// NewRegion returns new Region with xy
func NewRegion(x, y int) *Region {
	return &Region{
		X:          x,
		Y:          y,
		Data:       make([]byte, InformationSector),
		Locations:  make([]*Location, ChunkCount),
		Timestamps: make([]int32, ChunkCount),
	}
}

// Region is a section had 32x32 chunks
type Region struct {
	X int
	Y int

	Data []byte

	Locations  []*Location
	Timestamps []int32
}

func (Region) vaild(x, y int) error {
	if x < 0 || x >= 32 || y < 0 || y >= 32 {
		return fmt.Errorf("invaild x and y, they should be been 0 <= x < 32")
	}

	return nil
}

// getIndex returns index for Locations and Timestamps
func (Region) getIndex(x, y int) int {
	return x + (y * 32)
}

// Load loads a region from region file bytes
func (reg *Region) Load(b []byte) error {
	if len(b) < InformationSector {
		return fmt.Errorf("level.anvil.region: the region bytes isn't not enough")
	}

	reg.Data = b

	stream := binary.NewStreamBytes(b[:LocationsBytes+TimestampsBytes])

	// Read Locations
	reg.Locations = make([]*Location, ChunkCount)

	for i := 0; i < len(reg.Locations); i++ {
		offset, err := stream.Triad() // location of chunk data
		if err != nil {
			return err
		}

		count, err := stream.Byte() // lenght of chunk data
		if err != nil {
			return err
		}

		reg.Locations[i] = &Location{
			Off:   offset,
			Count: count,
		}
	}

	// Read timestamps
	reg.Timestamps = make([]int32, ChunkCount)

	for i := 0; i < len(reg.Timestamps); i++ {
		stamp, err := stream.Int()
		if err != nil {
			return err
		}

		reg.Timestamps[i] = stamp
	}

	return nil
}

// Save saves region data, returns bytes for a region file
func (reg *Region) Save() ([]byte, error) {
	// TODO: write

	return nil, nil
}

// ReadChunk reads a chunk, returns chunk data as []byte
// If the chunk doesn't exist, returns nil both two values
func (reg *Region) ReadChunk(x, y int) ([]byte, error) {
	err := reg.vaild(x, y)
	if err != nil {
		return nil, err
	}

	locat := reg.Locations[reg.getIndex(x, y)]

	off := int(locat.Off) * Sector
	ln := int(locat.Count) * Sector

	if off == 0 { // It haven't generated yet
		return nil, nil
	}

	stream := binary.NewStreamBytes(reg.Data[off : off+ln]) // chunk data + pads

	realLen, err := stream.Int() // ln = readLen + pad(4096 - (readlen % 4096)
	if err != nil {
		return nil, err
	}

	ctype, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	data := stream.Get(int(realLen)) // chunk data

	switch ctype {
	case CompressionGZip:
		read, err := gzip.NewReader(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		defer read.Close()

		data, err = ioutil.ReadAll(read)
		if err != nil {
			return nil, err
		}
	case CompressionZlib:
		read, err := zlib.NewReader(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		defer read.Close()

		data, err = ioutil.ReadAll(read)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// Location is a location info for chunk data
type Location struct {
	Off   binary.Triad
	Count byte
}
