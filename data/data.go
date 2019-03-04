package data

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"os"

	"github.com/beito123/nbt"
)

// LoadData loads level.dat file, returns as *nbt.Compound
func LoadData(path string) (*nbt.Compound, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	
	defer file.Close()

	return nil, nil
}
