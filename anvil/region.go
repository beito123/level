package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"

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
	path := rl.toRegionFile(x, y)

	if !util.ExistFile(path) {
		if !create {
			return nil, errors.New("level.anvil: couldn't find the region")
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

	reader := bufio.NewReader(file)

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n == 0 {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	err = reg.Load(buf)
	if err != nil {
		return nil, err
	}

	return reg, nil
}

// Region is a section had 32x32 chunks
type Region struct {
	X int
	Y int
}

// Load loads a region from region file bytes
func (reg *Region) Load(b []byte) error {
	// TODO: write

	return nil
}
