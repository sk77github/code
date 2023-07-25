walkdir_get_files:
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

func main() {
	var SystemDir string

	flag.StringVar(&SystemDir, "dir", "", " give the dir for the root ")
	flag.Parse()
	if SystemDir == "" {

		panic("SystemDir is empty")
		return
	}

	files, _ := ioutil.ReadDir(SystemDir)
	for _, file := range files {
		if file.IsDir() {

			dir := filepath.Join(SystemDir, file.Name())
			filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				fmt.Println(path)
				return nil
			})
		}
	}

}
