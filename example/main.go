package main

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/beito123/level/util"

	"github.com/beito123/level/anvil"
)

func main() {
	// Test :P

	err := test()
	if err != nil {
		panic(err)
	}
}

func test() error {
	resPath := "./resources/Vanilla_Resource_Pack_1.9.0"

	bbytes, err := ioutil.ReadFile(resPath + "/blocks.json")
	if err != nil {
		return err
	}

	regCommentLine := regexp.MustCompile(`//.*\n`)
	bbytes = regCommentLine.ReplaceAll(bbytes, []byte{}) // Remove

	var data map[string]interface{}

	err = json.Unmarshal(bbytes, &data)
	if err != nil {
		return err
	}

	blockList := make(map[string]string)

	for name, d := range data {
		var tname string

		d2, ok := d.(map[string]interface{})
		if !ok {
			continue // ignore // it's format_version
		}

		switch ntype := d2["textures"].(type) {
		case map[string]interface{}:
			tname = ntype["up"].(string)
		case string:
			tname = ntype
		default:
			fmt.Printf("unknown: %#v\n", ntype)
			continue
		}

		blockList["minecraft:"+name] = resPath + "/textures/blocks/" + tname + ".png"
	}

	loader, err := anvil.NewRegionLoader("./region-v1.13", anvil.RegionFileAnvil)
	if err != nil {
		return err
	}

	region, err := loader.LoadRegion(0, 0, false)
	if err != nil {
		return err
	}

	b, err := region.ReadChunk(0, 0)
	if err != nil {
		return err
	}

	chunk, err := anvil.ReadChunk(0, 0, b)
	if err != nil {
		return err
	}

	//targetY := 8

	subchunks := chunk.SubChunks()

	maker := ChunkImageMaker{}
	maker.Ready()

	for _, sub := range subchunks {
		if sub == nil {
			continue
		}

		for i, bs := range sub.Palette {
			tpath, ok := blockList[bs.Name]
			if !ok {
				continue
			}

			if !util.ExistFile(tpath) {
				continue // ignore
			}

			file, err := os.Open(tpath)
			if err != nil {
				return err
			}

			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				return err
			}

			maker.AddBlockData(i, img)

			fmt.Printf("Added Palette(id:%d, name:%s)\n", i, bs.Name)
		}

		fmt.Printf("Subchunks Y:%d\n", sub.Y)
		for y := 0; y < 16; y++ {
			/*if y != targetY {
				continue
			}*/

			fmt.Printf("block: y:%d\n", int(16*sub.Y)+y)

			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {
					fmt.Printf("%d,", sub.Blocks[y<<8|z<<4|x])

					maker.Add(x, z, int(sub.Blocks[y<<8|z<<4|x]))
				}

				fmt.Printf("\n")
			}
		}

	}
	maker.Output(fmt.Sprintf("./chunk.png"))

	fmt.Printf("jagajaga")

	return nil
}

type ChunkImageMaker struct {
	Image *image.RGBA

	BlockList map[int]image.Image
}

func (mk *ChunkImageMaker) Ready() {
	line := 16 * 16
	mk.Image = image.NewRGBA(image.Rect(0, 0, line, line))
	mk.BlockList = make(map[int]image.Image)
}

func (mk *ChunkImageMaker) Output(path string) error {
	file, _ := os.Create(path)
	defer file.Close()

	return png.Encode(file, mk.Image)
}

func (mk *ChunkImageMaker) Add(x, y, id int) {
	block, ok := mk.BlockList[id]
	if !ok {
		return
	}

	SetImage(block, mk.Image, x*16, y*16)

	return
}

func (mk *ChunkImageMaker) AddBlockData(id int, img image.Image) {
	mk.BlockList[id] = img
}

type BlockData struct {
	Textures *interface{} `json:"textures"`
}

func SetImage(src image.Image, dst *image.RGBA, atX, atY int) {
	for y := 0; y < src.Bounds().Dy(); y++ { // y
		for x := 0; x < src.Bounds().Dx(); x++ { // x
			dst.Set(atX+x, atY+y, src.At(x, y))
		}
	}
}
