package block

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"strconv"
)

// List is a compatibility data for block
type List map[string]*Block

// BlockListV112 is a block data for v1.12
var BlockListV112 = LoadV112()

// BlockListV113 is a block data for v1.13
var BlockListV113 = LoadV113()

// Block is a common block data
// TODO: support level.BlockState
type Block struct {
	Name       string // minecraft:dirt
	Properties map[string]string

	ID   int
	Meta int
}

// FromBlockID returns Block from old block id and meta
func FromBlockID(id int, meta int) *Block {
	data, ok := BlockListV112[ToNumberID(id)]
	if !ok {
		return nil
	}

	name := GetV112ToV113(data.Name, meta)

	return &Block{
		Name:       name,
		Properties: make(map[string]string), // TODO: support
	}
}

func ToNumberID(id int) string {
	return strconv.Itoa(id)
}

func ToNumberIDMeta(id int, meta int) string {
	return strconv.Itoa(id) + ":" + strconv.Itoa(meta)
}

// MinecraftPrefix is a block prefix
const MinecraftPrefix = "minecraft:"

// GetV112ToV113 gets a block name which is converted from v1.12 to v1.13
func GetV112ToV113(name string, meta int) string {
	data, ok := V112ToV113[name+":"+strconv.Itoa(meta)]
	if !ok {
		data, ok := V112ToV113[name]
		if !ok {
			return name
		}

		return data
	}

	return data
}
