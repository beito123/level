package asset

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"io"
	"net/http"
	"os"
)

func OpenResource(path string) (http.File, error) {
	return Assets.Open(path)
}

func ExistResource(path string) bool {
	_, err := OpenResource(path)
	if err != nil {
		return false
	}

	return true
}

func CopyResource(path string, to string) error {
	src, err := OpenResource(path)
	if err != nil {
		return err
	}

	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
