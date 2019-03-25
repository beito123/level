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
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
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

var making *MakingImage

type MakingImage struct {
	Image  *gif.GIF
	Bounds image.Rectangle
	Point  image.Point

	Delay int
}

func (mak *MakingImage) Ready() {
	mak.Image = &gif.GIF{}
}

func (mak *MakingImage) Outputs(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	err = gif.EncodeAll(file, mak.Image)
	if err != nil {
		return err
	}

	return nil
}

func (mak *MakingImage) Add(img *image.RGBA, index int) {
	/*
		pimg := image.NewPaletted(mak.Bounds, palette.Plan9)
		pt := mak.Point
		bounds := img.Bounds()
		rect := image.Rect(pt.X, pt.Y, pt.X+bounds.Dx(), pt.Y+bounds.Dy())
		draw.Draw(pimg, rect, img, image.ZP, draw.Over)
	*/

	pt := mak.Point
	bounds := img.Bounds()
	rect := image.Rect(pt.X, pt.Y, pt.X+bounds.Dx(), pt.Y+bounds.Dy())
	if index >= len(mak.Image.Image) { // New
		pimg := image.NewPaletted(mak.Bounds, palette.WebSafe)
		draw.Draw(pimg, rect, img, image.ZP, draw.Src)

		mak.Image.Delay = append(mak.Image.Delay, mak.Delay)
		mak.Image.Image = append(mak.Image.Image, pimg)
	} else { // Over
		pimg := mak.Image.Image[index]
		draw.Draw(pimg, rect, img, image.ZP, draw.Over)
	}
}

func test() error {
	resPath := "./resources/Vanilla_Resource_Pack_1.9.0"
	regionPath := "./region-v1.13"

	generator, err := NewMapGenerator(resPath, regionPath)
	if err != nil {
		return err
	}

	generator.Textures.AddAlias("minecraft:air", "minecraft:cave_air")

	generator.Textures.PathList["minecraft:granite"] = resPath + "/textures/blocks/" + "stone_granite.png"
	generator.Textures.PathList["minecraft:diorite"] = resPath + "/textures/blocks/" + "stone_diorite.png"
	generator.Textures.PathList["minecraft:andesite"] = resPath + "/textures/blocks/" + "stone_andesite.png"
	generator.Textures.PathList["minecraft:lava"] = resPath + "/textures/blocks/" + "lava_placeholder.png"
	generator.Textures.PathList["minecraft:water"] = resPath + "/textures/blocks/" + "water_placeholder.png"
	generator.Textures.PathList["minecraft:grass"] = resPath + "/textures/blocks/" + "grass_carried.png"

	scale := 4
	line := 16 * 16 * scale
	img := image.NewRGBA(image.Rect(0, 0, line, line))

	bx := -8
	by := -8
	//base := 0

	making = &MakingImage{
		Delay:  1,
		Bounds: image.Rect(0, 0, line, line),
	}

	making.Ready()

	//bar := pb.StartNew(scale * scale)
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			x := bx + i
			y := by + j

			making.Point = image.Pt(i*16*16, j*16*16)

			gimg, err := generator.Generate(x, y)
			if err != nil {
				return err
			}

			//bar.Increment()

			if gimg == nil {
				continue
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

	err = making.Outputs("./ani.gif")
	if err != nil {
		return err
	}

	return nil
}

// NewMapGenerator returns new MapGenerator
// path is a dir path for offical resource pack
// rpath is a region dir path
func NewMapGenerator(path, rpath string) (*MapGenerator, error) {
	tm := NewTextureManager()

	err := tm.LoadResourcePack(path)
	if err != nil {
		return nil, err
	}

	return &MapGenerator{
		Path:     rpath,
		Textures: tm,
	}, nil
}

type MapGenerator struct {
	Path     string
	Textures *TextureManager

	loader *anvil.RegionLoader
}

func (mg *MapGenerator) HasLoaded() bool {
	return mg.loader != nil
}

func (mg *MapGenerator) LoadRegions() error {
	loader, err := anvil.NewRegionLoader(mg.Path, anvil.RegionFileAnvil)
	if err != nil {
		return err
	}

	mg.loader = loader

	return nil
}

// Generate generates a chunk image
// path is a dir path for region
// x and y are chunk coordinates
// if it's returned nil as Image, the chunk is not created
func (mg *MapGenerator) Generate(x, y int) (image.Image, error) {
	if !mg.HasLoaded() {
		err := mg.LoadRegions()
		if err != nil {
			return nil, err
		}
	}

	region, err := mg.loader.LoadRegion(x>>5, y>>5, false)
	if err != nil {
		return nil, err
	}

	b, err := region.ReadChunk(x&31, y&31)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	chunk, err := anvil.ReadChunk(x&15, y&15, b)
	if err != nil {
		return nil, err
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
			if !mg.Textures.HasTexture(bs.Name) {
				//fmt.Printf("Ignore palette(id:%d, name: %s)\n", i, bs.Name)
				continue
			}

			img, err := mg.Textures.GetTexture(bs.Name)
			if err != nil {
				return nil, fmt.Errorf("happened errors while processing palette(id:%d, name: %s) error:%s", i, bs.Name, err)
			}

			//fmt.Printf("Added palette(id:%d, name: %s)\n", i, bs.Name)

			maker.AddBlockData(i, img)
		}

		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {
					maker.Add(x, z, int(sub.Blocks[y<<8|z<<4|x]))

					id := int(sub.Blocks[y<<8|z<<4|x])
					if id >= len(sub.Palette) {
						fmt.Printf("invail palette, id: %d (0b%b)\n", id, id)
					}
				}
			}

			making.Add(maker.Image, (int(sub.Y)*16)+y)
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
}

func (tm *TextureManager) getBlockName(name string) (string, bool) {
	_, ok := tm.PathList[name]
	if ok {
		return name, true
	}

	for n, v := range tm.Aliases {
		for _, c := range v {
			if c == name {
				return n, true
			}
		}
	}

	return "", false
}

func (tm *TextureManager) AddAlias(name string, aliases ...string) bool {
	tm.Aliases[name] = append(tm.Aliases[name], aliases...)

	return true
}

func (tm *TextureManager) HasTexture(name string) bool {
	name, ok := tm.getBlockName(name)
	if !ok {
		return false
	}

	path, ok := tm.PathList[name]
	if !ok {
		return false
	}

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

	return tm.preparedImages[name], nil
}

func (tm *TextureManager) HasPrepared(name string) bool {
	name, ok := tm.getBlockName(name)
	if !ok {
		return false
	}

	_, ok = tm.preparedImages[name]

	return ok
}

func (tm *TextureManager) Prepare(name string) error {
	name, ok := tm.getBlockName(name)
	if !ok {
		return fmt.Errorf("couldn't find a block")
	}

	path, ok := tm.PathList[name]
	if !ok {
		return fmt.Errorf("couldn't find a path for the block")
	}

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

	tm.preparedImages[name] = img

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

func SetImage(src image.Image, dst *image.RGBA, atX, atY int) {
	for y := 0; y < src.Bounds().Dy(); y++ { // y
		for x := 0; x < src.Bounds().Dx(); x++ { // x
			dst.Set(atX+x, atY+y, src.At(x, y))
		}
	}
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
