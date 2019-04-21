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
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/beito123/level"
	"github.com/pkg/errors"

	"github.com/beito123/level/block"

	"github.com/beito123/level/leveldb"

	"github.com/pbnjay/pixfont"

	"github.com/beito123/level/util"
)

func main() {
	err := test()
	if err != nil {
		fmt.Printf("Error: %s", errors.WithStack(err))
	}
}

func test() error {
	resPath := "./resources/vanilla"

	lvl, err := leveldb.New("./db")
	if err != nil {
		return err
	}

	lvl.Format = &leveldb.ChunkFormatV120{}

	generator, err := NewMapGenerator(resPath, lvl)
	if err != nil {
		return err
	}

	generator.Textures.AddAlias("minecraft:air", "minecraft:cave_air")
	generator.Textures.AddAlias("minecraft:grass_block", "minecraft:grass")

	generator.Textures.PathList["minecraft:granite"] = resPath + "/textures/blocks/" + "stone_granite.png"
	generator.Textures.PathList["minecraft:diorite"] = resPath + "/textures/blocks/" + "stone_diorite.png"
	generator.Textures.PathList["minecraft:andesite"] = resPath + "/textures/blocks/" + "stone_andesite.png"
	generator.Textures.PathList["minecraft:lava"] = resPath + "/textures/blocks/" + "lava_placeholder.png"
	generator.Textures.PathList["minecraft:water"] = resPath + "/textures/blocks/" + "water_placeholder.png"
	generator.Textures.PathList["minecraft:grass"] = resPath + "/textures/blocks/" + "grass_carried.png"
	generator.Textures.PathList["minecraft:grass_block"] = resPath + "/textures/blocks/" + "grass_carried.png"

	scale := 8
	line := 16 * 16 * scale
	img := image.NewRGBA(image.Rect(0, 0, line, line))

	bx := 0
	by := 0
	//base := 0

	//bar := pb.StartNew(scale * scale)
	/*making := &MakingImage{
		Delay:  1,
		Bounds: image.Rect(0, 0, line, line),
	}

	making.Ready()*/

	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			x := bx + i
			y := by + j

			//making.Point = image.Pt(i*16*16, j*16*16)

			gimg, err := generator.Generate(x, y)
			if err != nil {
				return err
			}

			//bar.Increment()

			if gimg == nil {
				continue
			}

			rimg, ok := gimg.(*image.RGBA)
			if ok {
				pixfont.DrawString(rimg, 8, 8, fmt.Sprintf("%d, %d", x, y), color.Black)
			}

			SetImage(gimg, img, i*16*16, j*16*16)
		}
	}

	//bar.FinishPrint("complete!")

	path := "./chunks.png"

	file, _ := os.Create(path)
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

// NewMapGenerator returns new MapGenerator
// path is a dir path for offical resource pack
// rpath is a region dir path
func NewMapGenerator(path string, lvl level.Format) (*MapGenerator, error) {
	tm := NewTextureManager()

	err := tm.LoadResourcePack(path)
	if err != nil {
		return nil, err
	}

	return &MapGenerator{
		Level:    lvl,
		Textures: tm,
	}, nil
}

type MapGenerator struct {
	Level         level.Format
	Textures      *TextureManager
	EnabledMaking bool
	Making        []image.Image
}

func (mg *MapGenerator) Clone(tm *TextureManager) *MapGenerator {
	return &MapGenerator{
		Textures:      tm,
		EnabledMaking: mg.EnabledMaking,
		Making:        mg.Making,
	}
}

// Generate generates a chunk image
// path is a dir path for region
// x and y are chunk coordinates
// if it's returned nil as Image, the chunk is not created
func (mg *MapGenerator) Generate(x, y int) (image.Image, error) {
	if mg.Level.HasGeneratedChunk(x, y) {
		return nil, nil
	}

	if !mg.Level.IsLoadedChunk(x, y) {
		err := mg.Level.LoadChunk(x, y)
		if err != nil {
			return nil, err
		}
	}

	chunk, ok := mg.Level.Chunk(x, y)
	if !ok {
		return nil, fmt.Errorf("unknown error")
	}

	maker := ChunkImageMaker{}
	//maker.EnabledFreeMap = true
	maker.Ready()

	skipped := make(map[string]bool)

	/*c, ok := chunk.(*leveldb.Chunk)
	if ok {
		fmt.Printf("Fin: %d", c.Finalization)
	}*/

	for y := 0; y < 256; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				bl, err := chunk.GetBlock(x, y, z)
				if err != nil {
					return nil, err
				}

				var name string

				b, ok := block.BlockListV112[bl.Name()]
				if ok {
					name = b.Name
				} else {
					name = bl.Name()
				}

				if name != "minecraft:air" {
					_, ok := skipped[name]
					if !ok && !maker.HasBlockData(name) {
						if !mg.Textures.HasTexture(name) {
							skipped[name] = true
							fmt.Printf("Ignore palette(name: %s)\n", name)
							continue
						}

						img, err := mg.Textures.GetTexture(name)
						if err != nil {
							return nil, fmt.Errorf("happened errors while processing palette(name: %s) error:%s", name, err)
						}

						maker.AddBlockData(name, img)
					}

					maker.Add(x, z, name)
				}

				/*
					bl, err := sub.AtBlock(x, 15-y, z)
					if err != nil {
						return nil, err
					}

					//fmt.Printf("test2:%s\n", bl.ToBlockData().Name)
					name := bl.ToBlockData().Name
					if maker.IsFree(x, z) && name != "minecraft:air" {
						maker.Add(x, z, name)
					}

				*/
			}
		}
	}

	return maker.Image, nil
}

var regCommentLine = regexp.MustCompile(`//.*\n`)

func NewTextureManager() *TextureManager {
	return &TextureManager{
		PathList:       make(map[string]string),
		Aliases:        make(map[string][]string),
		preparedImages: make(map[string]image.Image),
	}
}

// TextureManager control textures for blocks
type TextureManager struct {
	PathList map[string]string
	Aliases  map[string][]string

	preparedImages map[string]image.Image

	mutex sync.RWMutex
}

func (tm *TextureManager) getBlockName(name string) (string, bool) {
	tm.mutex.RLock()
	_, ok := tm.PathList[name]
	if ok {
		tm.mutex.RUnlock()
		return name, true
	}

	for n, v := range tm.Aliases {
		for _, c := range v {
			if c == name {
				tm.mutex.RUnlock()
				return n, true
			}
		}
	}

	tm.mutex.RUnlock()

	return "", false
}

func (tm *TextureManager) AddAlias(name string, aliases ...string) bool {
	tm.mutex.Lock()
	tm.Aliases[name] = append(tm.Aliases[name], aliases...)
	tm.mutex.Unlock()

	return true
}

func (tm *TextureManager) HasTexture(name string) bool {
	name, ok := tm.getBlockName(name)
	if !ok {
		return false
	}

	tm.mutex.RLock()

	path, ok := tm.PathList[name]
	if !ok {
		tm.mutex.RUnlock()
		return false
	}

	tm.mutex.RUnlock()

	if !util.ExistFile(path) {
		return false
	}

	return true
}

func (tm *TextureManager) GetTexture(name string) (image.Image, error) {
	if !tm.HasPrepared(name) {
		if !tm.HasTexture(name) {
			return nil, fmt.Errorf("couldn't find a image file")
		}

		err := tm.Prepare(name)
		if err != nil {
			return nil, err
		}
	}

	tm.mutex.RLock()
	result := tm.preparedImages[name]
	tm.mutex.RUnlock()

	return result, nil
}

func (tm *TextureManager) HasPrepared(name string) bool {
	name, ok := tm.getBlockName(name)
	if !ok {
		return false
	}
	tm.mutex.RLock()
	_, ok = tm.preparedImages[name]
	tm.mutex.RUnlock()

	return ok
}

func (tm *TextureManager) Prepare(name string) error {
	name, ok := tm.getBlockName(name)
	if !ok {
		return fmt.Errorf("couldn't find a block")
	}

	tm.mutex.RLock()
	path, ok := tm.PathList[name]
	if !ok {
		tm.mutex.RUnlock()
		return fmt.Errorf("couldn't find a path for the block")
	}

	tm.mutex.RUnlock()

	if !util.ExistFile(path) {
		return fmt.Errorf("couldn't find a image file")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	tm.mutex.Lock()
	tm.preparedImages[name] = img
	tm.mutex.Unlock()

	return nil
}

// LoadResourcePack loads textures from offical resource pack (you can download from https://www.minecraft.net/en-us/)
// path is a path for resource pack, you need to unzip in advance
func (tm *TextureManager) LoadResourcePack(path string) error {
	path = filepath.Clean(path)

	b, err := ioutil.ReadFile(path + "/blocks.json")
	if err != nil {
		return err
	}

	// bad hack for mojang // json isn't allowed comment lines (//)
	b = regCommentLine.ReplaceAll(b, []byte{}) // Remove comment lines

	// bad hack for mojang // you don't change type(string, array) for the same "textures"
	// I want to make struct...
	var data map[string]interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	tm.mutex.Lock()
	for name, d := range data { // of course, bad hack for mojang
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
			fmt.Printf("unknown: %s -> %#v\n", name, ntype)
			continue
		}

		tm.PathList["minecraft:"+name] = util.To(path, "/textures/blocks/"+tname+".png")
	}

	tm.mutex.Unlock()

	return nil
}

type ChunkImageMaker struct {
	Image *image.RGBA

	BlockList      map[string]image.Image
	FreeMap        []bool
	EnabledFreeMap bool
}

func (mk *ChunkImageMaker) Ready() {
	line := 16 * 16
	mk.Image = image.NewRGBA(image.Rect(0, 0, line, line))
	mk.BlockList = make(map[string]image.Image)

	mk.FreeMap = make([]bool, 16*16)
}

func (mk *ChunkImageMaker) Output(path string) error {
	file, _ := os.Create(path)
	defer file.Close()

	return png.Encode(file, mk.Image)
}

func (mk *ChunkImageMaker) IsFree(x, y int) bool {
	return !mk.FreeMap[y<<4|x]
}

func (mk *ChunkImageMaker) IsFull() bool {
	for _, v := range mk.FreeMap {
		if !v {
			return false
		}
	}

	return true
}

func (mk *ChunkImageMaker) Add(x, y int, name string) {
	block, ok := mk.BlockList[name]
	if !ok {
		return
	}

	if mk.EnabledFreeMap {
		mk.FreeMap[y<<4|x] = true
	}

	SetImage(block, mk.Image, x*16, y*16)

	return
}

func (mk *ChunkImageMaker) ResetBlockData() {
	mk.BlockList = make(map[string]image.Image)
}

func (mk *ChunkImageMaker) HasBlockData(name string) bool {
	_, ok := mk.BlockList[name]

	return ok
}

func (mk *ChunkImageMaker) AddBlockData(name string, img image.Image) {
	mk.BlockList[name] = img
}

func SetImage(src image.Image, dst *image.RGBA, atX, atY int) {
	/*for y := 0; y < src.Bounds().Dy(); y++ { // y
		for x := 0; x < src.Bounds().Dx(); x++ { // x
			dst.Set(atX+x, atY+y, src.At(x, y))
		}
	}*/

	size := src.Bounds().Size()
	rect := image.Rect(atX, atY, atX+size.X, atY+size.Y)
	draw.Draw(dst, rect, src, image.ZP, draw.Over)
}

/* Old codes

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

	subchunks := chunk.SubChunks()

	maker := ChunkImageMaker{}
	maker.Ready()

	for _, sub := range subchunks {
		if sub == nil {
			continue
		}

		maker.ResetBlockData()

		for i, bs := range sub.Palette {
			tpath, ok := blockList[bs.Name]
			if !ok {
				fmt.Printf("Ignore Palette(id:%d, name: %s), no listed\n", i, bs.Name)
				continue
			}

			if !util.ExistFile(tpath) {
				fmt.Printf("Ignore Palette(id:%d, name:%s), not found\n", i, bs.Name)
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

func (mk *ChunkImageMaker) ResetBlockData() {
	mk.BlockList = make(map[int]image.Image)
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
}*/
