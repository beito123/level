package util

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"math"
	"os"
	"path/filepath"
)

// File

func GetDir(path string) string {
	return filepath.Dir(filepath.Clean(path))
}

func ExistFile(file string) bool {
	f, err := os.Stat(file)
	return err == nil && !f.IsDir()
}

func ExistDir(dir string) bool {
	f, err := os.Stat(dir)

	return err == nil && f.IsDir()
}

func To(root string, child string) string {
	return root + string(filepath.Separator) + child
}

// Math

func CeilInt(x float64) int {
	return int(math.Ceil(x))
}
