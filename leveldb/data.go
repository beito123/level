package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/beito123/binary"
	"github.com/beito123/level"
	"github.com/beito123/nbt"
)

// DefaultProperties is the default properties for level.dat
var DefaultProperties = &Properties{
	Data: &nbt.Compound{
		Value: map[string]nbt.Tag{
			TagLevelName: nbt.NewStringTag("", ""),
			TagGameType:  nbt.NewByteTag("", int8(level.Survival)),
			TagSpawnX:    nbt.NewIntTag("", 0),
			TagSpawnY:    nbt.NewIntTag("", 0),
			TagSpawnZ:    nbt.NewIntTag("", 0),
		},
	},
	Version: 8,
}

var (
	TagLevelName = "LevelName"
	TagGameType  = "GameType"
	TagSpawnX    = "SpawnX"
	TagSpawnY    = "SpawnY"
	TagSpawnZ    = "SpawnZ"
)

// LoadLevelData loads properties from level.dat
func LoadLevelData(path string) (*Properties, error) {
	buf, err := ioutil.ReadFile(filepath.Join(path, LevelDataFile))
	if err != nil {
		return nil, err
	}

	stream := binary.NewOrderStreamBytes(binary.LittleEndian, buf)

	ver, err := stream.Int()
	if err != nil {
		return nil, err
	}

	ln, err := stream.Int()
	if err != nil {
		return nil, err
	}

	nstream := nbt.NewStreamBytes(nbt.LittleEndian, stream.Get(int(ln)))

	tag, err := nstream.ReadTag()
	if err != nil {
		return nil, err
	}

	com, ok := tag.(*nbt.Compound)
	if !ok {
		return nil, errors.New("unexpected " + nbt.GetTagName(tag.ID()) + "Tag, expected CompoundTag")
	}

	return &Properties{
		Data:    com,
		Version: int(ver),
	}, nil
}

// SaveLevelData saves properties to level.dat
func SaveLevelData(path string, pro *Properties) error {
	nstream := nbt.NewStream(nbt.LittleEndian)

	err := nstream.WriteTag(pro.Data)
	if err != nil {
		return err
	}

	buf := nstream.Bytes()
	ln := len(buf)

	stream := binary.NewOrderStream(nbt.LittleEndian)

	err = stream.PutInt(int32(pro.Version)) // uint?
	if err != nil {
		return err
	}

	err = stream.PutInt(int32(ln)) // uint?
	if err != nil {
		return err
	}

	err = stream.Put(buf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(path, LevelDataFile), stream.AllBytes(), os.ModePerm)
}

// Properties is data of level from level.dat
type Properties struct {
	Data    *nbt.Compound
	Version int
}
