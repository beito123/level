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

func setPrefix(prefix string, m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[prefix+k] = prefix + v
	}

	return n
}

// V112ToV113 is compatibility data between v1.12 and v1.13 with minecraft prefix
var V112ToV113 = setPrefix(MinecraftPrefix, v112Tov113)

// v112Tov113 is compatibility data between v1.12 and v1.13
var v112Tov113 = map[string]string{
	/*
		* TODO
		flower_pot
		skull
		standing_banner
		banner
		wall_banner
		*
		"infested_stone":                 []string{"stone_monster_egg"},
		"infested_cobblestone":           []string{"cobblestone_monster_egg"},
		"infested_stone_bricks":          []string{"stone_brick_monster_egg"},
		"infested_cracked_stone_bricks":  []string{"mossy_stone_brick_monster_egg"},
		"infested_mossy_stone_bricks":    []string{"cracked_stone_brick_monster_egg"},
		"infested_chiseled_stone_bricks": []string{"chiseled_stone_brick_monster_egg"},
		"skeleton_wall_skull":            []string{"skeleton_skull"},
		"zombie_wall_head":               []string{"zombie_skull"},
		"player_wall_head":               []string{"player_skull"},
		"creeper_wall_head":              []string{"creeper_skull"},
		"dragon_wall_head":               []string{"dragon_skull"},
		"white_banner":                   []string{"standing_banner"},
		"orange_banner":                  []string{"standing_banner"},
		"magenta_banner":                 []string{"standing_banner"},
		"light_blue_banner":              []string{"standing_banner"},
		"yellow_banner":                  []string{"standing_banner"},
		"lime_banner":                    []string{"standing_banner"},
		"pink_banner":                    []string{"standing_banner"},
		"gray_banner":                    []string{"standing_banner"},
		"light_gray_banner":              []string{"standing_banner"},
		"cyan_banner":                    []string{"standing_banner"},
		"purple_banner":                  []string{"standing_banner"},
		"blue_banner":                    []string{"standing_banner"},
		"brown_banner":                   []string{"standing_banner"},
		"green_banner":                   []string{"standing_banner"},
		"red_banner":                     []string{"standing_banner"},
		"black_banner":                   []string{"standing_banner"},
		"white_wall_banner":              []string{"wall_banner"},
		"orange_wall_banner":             []string{"wall_banner"},
		"magenta_wall_banner":            []string{"wall_banner"},
		"light_blue_wall_banner":         []string{"wall_banner"},
		"yellow_wall_banner":             []string{"wall_banner"},
		"lime_wall_banner":               []string{"wall_banner"},
		"pink_wall_banner":               []string{"wall_banner"},
		"gray_wall_banner":               []string{"wall_banner"},
		"light_gray_wall_banner":         []string{"wall_banner"},
		"cyan_wall_banner":               []string{"wall_banner"},
		"purple_wall_banner":             []string{"wall_banner"},
		"blue_wall_banner":               []string{"wall_banner"},
		"brown_wall_banner":              []string{"wall_banner"},
		"green_wall_banner":              []string{"wall_banner"},
		"red_wall_banner":                []string{"wall_banner"},
		"black_wall_banner":              []string{"wall_banner"}
	*/

	// Stone
	"stone:0": "stone",
	"stone:1": "granite",
	"stone:2": "polished_granite",
	"stone:3": "diorite",
	"stone:4": "polished_diorite",
	"stone:5": "andesite",
	"stone:6": "polished_andesite",

	// Grass block
	"grass": "grass_block",

	// Dirt
	"dirt:0": "dirt",
	"dirt:1": "coarse_dirt",
	"dirt:2": "podzol",

	// WoodenPlanks
	"planks:0": "oak_planks",
	"planks:1": "spruce_planks",
	"planks:2": "birch_planks",
	"planks:3": "jungle_planks",
	"planks:4": "acacia_planks",
	"planks:5": "dark_oak_planks",

	// Sapling
	"sapling:0": "oak_sapling",
	"sapling:1": "spruce_sapling",
	"sapling:2": "birch_sapling",
	"sapling:3": "jungle_sapling",
	"sapling:4": "acacia_sapling",
	"sapling:5": "dark_oak_sapling",

	// Sand
	"sand:0": "sand",
	"sand:1": "red_sand",

	// Log
	"log:0":   "oak_log",
	"log:4":   "oak_log",
	"log:8":   "oak_log",
	"log:12":  "oak_log",
	"log:1":   "spruce_log",
	"log:5":   "spruce_log",
	"log:9":   "spruce_log",
	"log:13":  "spruce_log",
	"log:2":   "birch_log",
	"log:6":   "birch_log",
	"log:10":  "birch_log",
	"log:14":  "birch_log",
	"log:3":   "jungle_log",
	"log:7":   "jungle_log",
	"log:11":  "jungle_log",
	"log:15":  "jungle_log",
	"log2:0":  "acacia_log",
	"log2:4":  "acacia_log",
	"log2:8":  "acacia_log",
	"log2:12": "acacia_log",
	"log2:1":  "dark_oak_log",
	"log2:5":  "dark_oak_log",
	"log2:9":  "dark_oak_log",
	"log2:13": "dark_oak_log",

	// Leaves
	"leaves:0":   "oak_leaves",
	"leaves:4":   "oak_leaves",
	"leaves:8":   "oak_leaves",
	"leaves:12":  "oak_leaves",
	"leaves:1":   "spruce_leaves",
	"leaves:5":   "spruce_leaves",
	"leaves:9":   "spruce_leaves",
	"leaves:13":  "spruce_leaves",
	"leaves:2":   "birch_leaves",
	"leaves:6":   "birch_leaves",
	"leaves:10":  "birch_leaves",
	"leaves:14":  "birch_leaves",
	"leaves:3":   "jungle_leaves",
	"leaves:7":   "jungle_leaves",
	"leaves:11":  "jungle_leaves",
	"leaves:15":  "jungle_leaves",
	"leaves2:0":  "acacia_leaves",
	"leaves2:4":  "acacia_leaves",
	"leaves2:8":  "acacia_leaves",
	"leaves2:12": "acacia_leaves",
	"leaves2:1":  "dark_oak_leaves",
	"leaves2:5":  "dark_oak_leaves",
	"leaves2:9":  "dark_oak_leaves",
	"leaves2:13": "dark_oak_leaves",

	// Sponge
	"sponge:0": "sponge",
	"sponge:1": "wet_sponge",

	// Sandstone
	"sandstone:0": "sandstone",
	"sandstone:1": "chiseled_sandstone",
	"sandstone:2": "cut_sandstone",

	// NoteBlock
	"noteblock": "note_block",

	// PoweredRail
	"powered_rail": "golden_rail",

	// Cobweb
	"cobweb": "web",

	// Grass
	"tallgrass":   "dead_bush",
	"tallgrass:0": "dead_bush",
	"tallgrass:1": "grass",
	"tallgrass:2": "fern",
	"deadbush":    "dead_bush",

	// MoveingPiston
	"piston_extension": "moving_piston",

	// Wool
	"wool:0":  "white_wool",
	"wool:1":  "orange_wool",
	"wool:2":  "magenta_wool",
	"wool:3":  "light_blue_wool",
	"wool:4":  "yellow_wool",
	"wool:5":  "lime_wool",
	"wool:6":  "pink_wool",
	"wool:7":  "gray_wool",
	"wool:8":  "light_gray_wool",
	"wool:9":  "cyan_wool",
	"wool:10": "purple_wool",
	"wool:11": "blue_wool",
	"wool:12": "brown_wool",
	"wool:13": "green_wool",
	"wool:14": "red_wool",
	"wool:15": "black_wool",

	// Flower
	"yellow_flower": "dandelion",
	"red_flower:0":  "poppy",
	"red_flower:1":  "blue_orchid",
	"red_flower:2":  "allium",
	"red_flower:3":  "azure_bluet",
	"red_flower:4":  "red_tulip",
	"red_flower:5":  "orange_tulip",
	"red_flower:6":  "white_tulip",
	"red_flower:7":  "pink_tulip",
	"red_flower:8":  "oxeye_daisy",

	// Wooden slab

	"wooden_slab:0":  "oak_slab",
	"wooden_slab:8":  "oak_slab",
	"wooden_slab:1":  "spruce_slab",
	"wooden_slab:9":  "spruce_slab",
	"wooden_slab:2":  "birch_slab",
	"wooden_slab:10": "birch_slab",
	"wooden_slab:3":  "jungle_slab",
	"wooden_slab:11": "jungle_slab",
	"wooden_slab:4":  "acacia_slab",
	"wooden_slab:12": "acacia_slab",
	"wooden_slab:5":  "dark_oak_slab",
	"wooden_slab:13": "dark_oak_slab",

	// Stone slab
	"stone_slab:0":         "stone_slab",
	"stone_slab:8":         "stone_slab",
	"double_stone_slab:0":  "stone_slab",
	"stone_slab:1":         "sandstone_slab",
	"stone_slab:9":         "sandstone_slab",
	"double_stone_slab:1":  "sandstone_slab",
	"stone_slab:2":         "petrified_oak_slab",
	"stone_slab:10":        "petrified_oak_slab",
	"double_stone_slab:2":  "petrified_oak_slab",
	"stone_slab:3":         "cobblestone_slab",
	"stone_slab:11":        "cobblestone_slab",
	"double_stone_slab:3":  "cobblestone_slab",
	"stone_slab:4":         "brick_slab",
	"stone_slab:12":        "brick_slab",
	"double_stone_slab:4":  "brick_slab",
	"stone_slab:5":         "stone_brick_slab",
	"stone_slab:13":        "stone_brick_slab",
	"double_stone_slab:5":  "stone_brick_slab",
	"stone_slab:6":         "nether_brick_slab",
	"stone_slab:14":        "nether_brick_slab",
	"double_stone_slab:6":  "nether_brick_slab",
	"stone_slab:7":         "quartz_slab",
	"stone_slab:15":        "quartz_slab",
	"double_stone_slab:7":  "quartz_slab",
	"double_stone_slab:8":  "smooth_stone",
	"double_stone_slab:9":  "smooth_sandstone",
	"double_stone_slab:15": "smooth_quartz",
	"stone_slab2:0":        "red_sandstone_slab",
	"stone_slab2:8":        "red_sandstone_slab",
	"double_stone_slab2:0": "red_sandstone_slab",
	"double_stone_slab2:8": "smooth_red_sandstone",
	"purpur_slab":          "purpur_slab",
	"double_purpur_slab":   "purpur_slab",

	// Bricks block
	"bricks_block": "bricks",

	// Spawner
	"mob_spawner": "spawner",

	// Nether portal
	"portal": "nether_portal",

	// Torch
	"torch:0": "torch",
	"torch:1": "wall_torch",
	"torch:2": "wall_torch",
	"torch:3": "wall_torch",
	"torch:4": "wall_torch",

	// Furnace
	"lit_furnace": "furnace",

	// Stone stairs
	"cobblestone_stairs": "stone_stairs",

	// Wooden pressure plate
	"wooden_pressure_plate": "oak_pressure_plate",

	// Redstone ore
	"lit_redstone_ore": "redstone_ore",

	// Redstone torch
	"redstone_torch:0":       "redstone_wall_torch",
	"redstone_torch:1":       "redstone_wall_torch",
	"redstone_torch:2":       "redstone_wall_torch",
	"redstone_torch:3":       "redstone_wall_torch",
	"redstone_torch:4":       "redstone_torch",
	"unlit_redstone_torch:0": "redstone_wall_torch",
	"unlit_redstone_torch:1": "redstone_wall_torch",
	"unlit_redstone_torch:2": "redstone_wall_torch",
	"unlit_redstone_torch:3": "redstone_wall_torch",
	"unlit_redstone_torch:4": "redstone_torch",

	// Snow
	"snow_layer": "snow",
	"snow":       "snow_block",

	// Fence
	"fence": "oak_fence",

	// Pumpkin
	"pumpkin":     "carved_pumpkin",
	"lit_pumpkin": "jack_o_lantern",

	// Trapdoor
	"trapdoor": "oak_trapdoor",

	"monster_egg:0": "infested_stone",
	"monster_egg:1": "infested_cobblestone",
	"monster_egg:2": "infested_stone_bricks",
	"monster_egg:3": "infested_cracked_stone_bricks",
	"monster_egg:4": "infested_mossy_stone_bricks",
	"monster_egg:5": "infested_chiseled_stone_bricks",

	"stonebrick:0": "stone_bricks",
	"stonebrick:1": "mossy_stone_bricks",
	"stonebrick:2": "cracked_stone_bricks",
	"stonebrick:3": "chiseled_stone_bricks",

	"brown_mushroom_block:0":     "brown_mushroom_block",
	"brown_mushroom_block:1":     "brown_mushroom_block",
	"brown_mushroom_block:2":     "brown_mushroom_block",
	"brown_mushroom_block:3":     "brown_mushroom_block",
	"brown_mushroom_block:4":     "brown_mushroom_block",
	"brown_mushroom_block:5":     "brown_mushroom_block",
	"brown_mushroom_block:6":     "brown_mushroom_block",
	"brown_mushroom_block:7":     "brown_mushroom_block",
	"brown_mushroom_block:8":     "brown_mushroom_block",
	"brown_mushroom_block:9":     "brown_mushroom_block",
	"brown_mushroom_block:10":    "mushroom_stem",
	"brown_mushroom_block:11":    "brown_mushroom_block",
	"brown_mushroom_block:12":    "brown_mushroom_block",
	"brown_mushroom_block:13":    "brown_mushroom_block",
	"brown_mushroom_block:14":    "brown_mushroom_block",
	"brown_mushroom_block:15":    "mushroom_stem",
	"red_mushroom_block:0":       "red_mushroom_block",
	"red_mushroom_block:1":       "red_mushroom_block",
	"red_mushroom_block:2":       "red_mushroom_block",
	"red_mushroom_block:3":       "red_mushroom_block",
	"red_mushroom_block:4":       "red_mushroom_block",
	"red_mushroom_block:5":       "red_mushroom_block",
	"red_mushroom_block:6":       "red_mushroom_block",
	"red_mushroom_block:7":       "red_mushroom_block",
	"red_mushroom_block:8":       "red_mushroom_block",
	"red_mushroom_block:9":       "red_mushroom_block",
	"red_mushroom_block:10":      "mushroom_stem",
	"red_mushroom_block:11":      "red_mushroom_block",
	"red_mushroom_block:12":      "red_mushroom_block",
	"red_mushroom_block:13":      "red_mushroom_block",
	"red_mushroom_block:14":      "red_mushroom_block",
	"red_mushroom_block:15":      "mushroom_stem",
	"melon_block":                "melon",
	"fence_gate":                 "oak_fence_gate",
	"waterlily":                  "lily_pad",
	"nether_brick":               "nether_bricks",
	"end_bricks":                 "end_stone_bricks",
	"lit_redstone_lamp":          "redstone_lamp",
	"cobblestone_wall:0":         "cobblestone_wall",
	"cobblestone_wall:1":         "mossy_cobblestone_wall",
	"wooden_button":              "oak_button",
	"anvil:0":                    "anvil",
	"anvil:1":                    "anvil",
	"anvil:2":                    "anvil",
	"anvil:3":                    "anvil",
	"anvil:4":                    "chipped_anvil",
	"anvil:5":                    "chipped_anvil",
	"anvil:6":                    "chipped_anvil",
	"anvil:7":                    "chipped_anvil",
	"anvil:8":                    "damaged_anvil",
	"anvil:9":                    "damaged_anvil",
	"anvil:10":                   "damaged_anvil",
	"anvil:11":                   "damaged_anvil",
	"daylight_detector":          "daylight_detector",
	"daylight_detector_inverted": "daylight_detector",
	"quartz_ore":                 "nether_quartz_ore",
	"quartz_block:0":             "quartz_block",
	"quartz_block:1":             "chiseled_quartz_block",
	"quartz_block:2":             "quartz_pillar",
	"quartz_block:3":             "quartz_pillar",
	"quartz_block:4":             "quartz_pillar",
	"stained_hardened_clay:0":    "white_terracotta",
	"stained_hardened_clay:1":    "orange_terracotta",
	"stained_hardened_clay:2":    "magenta_terracotta",
	"stained_hardened_clay:3":    "light_blue_terracotta",
	"stained_hardened_clay:4":    "yellow_terracotta",
	"stained_hardened_clay:5":    "lime_terracotta",
	"stained_hardened_clay:6":    "pink_terracotta",
	"stained_hardened_clay:7":    "gray_terracotta",
	"stained_hardened_clay:8":    "light_gray_terracotta",
	"stained_hardened_clay:9":    "cyan_terracotta",
	"stained_hardened_clay:10":   "purple_terracotta",
	"stained_hardened_clay:11":   "blue_terracotta",
	"stained_hardened_clay:12":   "brown_terracotta",
	"stained_hardened_clay:13":   "green_terracotta",
	"stained_hardened_clay:14":   "red_terracotta",
	"stained_hardened_clay:15":   "black_terracotta",
	"carpet:0":                   "white_carpet",
	"carpet:1":                   "orange_carpet",
	"carpet:2":                   "magenta_carpet",
	"carpet:3":                   "light_blue_carpet",
	"carpet:4":                   "yellow_carpet",
	"carpet:5":                   "lime_carpet",
	"carpet:6":                   "pink_carpet",
	"carpet:7":                   "gray_carpet",
	"carpet:8":                   "light_gray_carpet",
	"carpet:9":                   "cyan_carpet",
	"carpet:10":                  "purple_carpet",
	"carpet:11":                  "blue_carpet",
	"carpet:12":                  "brown_carpet",
	"carpet:13":                  "green_carpet",
	"carpet:14":                  "red_carpet",
	"carpet:15":                  "black_carpet",
	"stained_glass:0":            "white_stained_glass",
	"stained_glass:1":            "orange_stained_glass",
	"stained_glass:2":            "magenta_stained_glass",
	"stained_glass:3":            "light_blue_stained_glass",
	"stained_glass:4":            "yellow_stained_glass",
	"stained_glass:5":            "lime_stained_glass",
	"stained_glass:6":            "pink_stained_glass",
	"stained_glass:7":            "gray_stained_glass",
	"stained_glass:8":            "light_gray_stained_glass",
	"stained_glass:9":            "cyan_stained_glass",
	"stained_glass:10":           "purple_stained_glass",
	"stained_glass:11":           "blue_stained_glass",
	"stained_glass:12":           "brown_stained_glass",
	"stained_glass:13":           "green_stained_glass",
	"stained_glass:14":           "red_stained_glass",
	"stained_glass:15":           "black_stained_glass",
	"stained_glass_pane:0":       "white_stained_glass_pane",
	"stained_glass_pane:1":       "orange_stained_glass_pane",
	"stained_glass_pane:2":       "magenta_stained_glass_pane",
	"stained_glass_pane:3":       "light_blue_stained_glass_pane",
	"stained_glass_pane:4":       "yellow_stained_glass_pane",
	"stained_glass_pane:5":       "lime_stained_glass_pane",
	"stained_glass_pane:6":       "pink_stained_glass_pane",
	"stained_glass_pane:7":       "gray_stained_glass_pane",
	"stained_glass_pane:8":       "light_gray_stained_glass_pane",
	"stained_glass_pane:9":       "cyan_stained_glass_pane",
	"stained_glass_pane:10":      "purple_stained_glass_pane",
	"stained_glass_pane:11":      "blue_stained_glass_pane",
	"stained_glass_pane:12":      "brown_stained_glass_pane",
	"stained_glass_pane:13":      "green_stained_glass_pane",
	"stained_glass_pane:14":      "red_stained_glass_pane",
	"stained_glass_pane:15":      "black_stained_glass_pane",
	"concrete:0":                 "white_concrete",
	"concrete:1":                 "orange_concrete",
	"concrete:2":                 "magenta_concrete",
	"concrete:3":                 "light_blue_concrete",
	"concrete:4":                 "yellow_concrete",
	"concrete:5":                 "lime_concrete",
	"concrete:6":                 "pink_concrete",
	"concrete:7":                 "gray_concrete",
	"concrete:8":                 "light_gray_concrete",
	"concrete:9":                 "cyan_concrete",
	"concrete:10":                "purple_concrete",
	"concrete:11":                "blue_concrete",
	"concrete:12":                "brown_concrete",
	"concrete:13":                "green_concrete",
	"concrete:14":                "red_concrete",
	"concrete:15":                "black_concrete",
	"concrete_powder:0":          "white_concrete_powder",
	"concrete_powder:1":          "orange_concrete_powder",
	"concrete_powder:2":          "magenta_concrete_powder",
	"concrete_powder:3":          "light_blue_concrete_powder",
	"concrete_powder:4":          "yellow_concrete_powder",
	"concrete_powder:5":          "lime_concrete_powder",
	"concrete_powder:6":          "pink_concrete_powder",
	"concrete_powder:7":          "gray_concrete_powder",
	"concrete_powder:8":          "light_gray_concrete_powder",
	"concrete_powder:9":          "cyan_concrete_powder",
	"concrete_powder:10":         "purple_concrete_powder",
	"concrete_powder:11":         "blue_concrete_powder",
	"concrete_powder:12":         "brown_concrete_powder",
	"concrete_powder:13":         "green_concrete_powder",
	"concrete_powder:14":         "red_concrete_powder",
	"concrete_powder:15":         "black_concrete_powder",
	"dye:0":                      "bone_meal",
	"dye:1":                      "orange_dye",
	"dye:2":                      "magenta_dye",
	"dye:3":                      "light_blue_dye",
	"dye:4":                      "dandelion_yellow",
	"dye:5":                      "lime_dye",
	"dye:6":                      "pink_dye",
	"dye:7":                      "gray_dye",
	"dye:8":                      "light_gray_dye",
	"dye:9":                      "cyan_dye",
	"dye:10":                     "purple_dye",
	"dye:11":                     "lapis_lazuli",
	"dye:12":                     "cocoa_beans",
	"dye:13":                     "cactus_green",
	"dye:14":                     "rose_red",
	"dye:15":                     "ink_sac",
	"bed:0":                      "white_bed",
	"bed:1":                      "orange_bed",
	"bed:2":                      "magenta_bed",
	"bed:3":                      "light_blue_bed",
	"bed:4":                      "yellow_bed",
	"bed:5":                      "lime_bed",
	"bed:6":                      "pink_bed",
	"bed:7":                      "gray_bed",
	"bed:8":                      "light_gray_bed",
	"bed:9":                      "cyan_bed",
	"bed:10":                     "purple_bed",
	"bed:11":                     "blue_bed",
	"bed:12":                     "brown_bed",
	"bed:13":                     "green_bed",
	"bed:14":                     "red_bed",
	"bed:15":                     "black_bed",
	"hardened_clay":              "terracotta",
	"slime":                      "slime_block",
	"double_plant:0":             "sunflower",
	"double_plant:1":             "lilac",
	"double_plant:2":             "tall_grass",
	"double_plant:3":             "large_fern",
	"double_plant:4":             "rose_bush",
	"double_plant:5":             "peony",
	"prismarine:1":               "prismarine",
	"prismarine:2":               "prismarine_bricks",
	"prismarine:3":               "dark_prismarine",
	"red_sandstone:0":            "red_sandstone",
	"red_sandstone:1":            "chiseled_red_sandstone",
	"red_sandstone:2":            "cut_red_sandstone",
	"magma":                      "magma_block",
	"red_nether_brick":           "red_nether_bricks",
	"silver_shulker_box":         "light_gray_shulker_box",
	"silver_glazed_terracotta":   "light_gray_glazed_terracotta",
	"wooden_door":                "oak_door",
	"powered_repeater":           "repeater",
	"unpowered_repeater":         "repeater",
	"powered_comparator":         "comparator",
	"unpowered_comparator":       "comparator",
	"coal:0":                     "coal",
	"coal:1":                     "charcoal",
	"golden_apple:0":             "golden_apple",
	"golden_apple:1":             "enchanted_golden_apple",
	"standing_sign":              "sign",
	"flowing_water":              "water",
	"flowing_lava":               "lava",
	"boat":                       "oak_boat",
	"reeds":                      "sugar_cane",
	"fish:0":                     "cod",
	"fish:1":                     "salmon",
	"fish:2":                     "tropical_fish",
	"fish:3":                     "pufferfish",
	"cooked_fish:0":              "cooked_cod",
	"cooked_fish:1":              "cooked_salmon",
	"melon":                      "melon_slice",
	"pumpkin_stem":               "pumpkin_stem",
	"melon_stem":                 "melon_stem",
	"speckled_melon":             "glistering_melon_slice",
	"fireworks":                  "firework_rocket",
	"firework_charge":            "firework_star",
	"netherbrick":                "nether_brick",
	"chorus_fruit_popped":        "popped_chorus_fruit",
	"record_13":                  "music_disc_13",
	"record_cat":                 "music_disc_cat",
	"record_blocks":              "music_disc_blocks",
	"record_chirp":               "music_disc_chirp",
	"record_far":                 "music_disc_far",
	"record_mall":                "music_disc_mall",
	"record_mellohi":             "music_disc_mellohi",
	"record_stal":                "music_disc_stal",
	"record_strad":               "music_disc_strad",
	"record_ward":                "music_disc_ward",
	"record_11":                  "music_disc_11",
	"record_wait":                "music_disc_wait",
}
