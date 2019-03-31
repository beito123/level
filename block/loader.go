package block

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/beito123/level/asset"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func handleError(err error) {
	panic(fmt.Errorf("level.block: happened errors while it's loading block data Error: %s", err.Error()))
}

// LoadV112 loads block data for v112
func LoadV112() List {
	file, err := asset.OpenResource("/static/v112/blocks.json")
	if err != nil {
		handleError(err)
	}

	defer file.Close()

	type Format struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		//DisplayName string `json:"displayName"`
		Variations []struct {
			Meta int `json:"metadata"`
			//DisplayName string `json:"displayName"`
		} `json:"variations"`
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		handleError(err)
	}

	var data []Format
	err = json.Unmarshal(b, &data)
	if err != nil {
		handleError(err)
	}

	list := make(List)
	for _, value := range data {
		name := MinecraftPrefix + value.Name

		list[name] = &BlockData{
			Name: name,
			//DisplayName: value.DisplayName,

			ID:   value.ID,
			Meta: 0,
		}

		list[ToNumberID(value.ID)] = list[name]

		for _, val := range value.Variations {
			n := name + ":" + strconv.Itoa(val.Meta)

			list[n] = &BlockData{
				Name: name,
				//DisplayName: val.DisplayName,
				ID:   value.ID,
				Meta: val.Meta,
			}

			nid := ToNumberIDMeta(value.ID, val.Meta)
			list[nid] = &BlockData{
				Name: name,
				//DisplayName: val.DisplayName,
				ID:   value.ID,
				Meta: val.Meta,
			}
		}
	}

	//
	/*

		for key, to := range cdata {
			val, ok := list[MinecraftPrefix+to]
			if !ok {
				continue
			}

			list[MinecraftPrefix+key] = val
		}*/

	return list
}

// LoadV113 loads block data for v113
func LoadV113() List {
	file, err := asset.OpenResource("/static/v113/blocks.json")
	if err != nil {
		handleError(err)
	}

	defer file.Close()

	type Format struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		handleError(err)
	}

	var data []Format
	err = json.Unmarshal(b, &data)
	if err != nil {
		handleError(err)
	}

	list := make(List)
	for _, value := range data {
		name := MinecraftPrefix + value.Name

		list[name] = &BlockData{
			Name: name,
			//DisplayName: value.DisplayName,

			ID:   value.ID,
			Meta: 0,
		}
	}

	/*
		// Debug
		debug, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Printf("test:\n%s", string(debug)):*/

	return list
}
