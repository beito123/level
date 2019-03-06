package main

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

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
	loader, err := anvil.NewRegionLoader("./region")
	if err != nil {
		return err
	}

	_, err = loader.LoadRegion(0, 0, false)
	if err != nil {
		return err
	}

	fmt.Printf("jagajaga")

	return nil
}
