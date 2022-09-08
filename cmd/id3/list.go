package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/blugnu/tags/mp3"
)

func list() {
	for _, arg := range flag.Args() {
		currdir := ""
		files := []string{}
		maxlen := 0
		filepath.Walk(arg, func(path string, info fs.FileInfo, err error) error {
			base := filepath.Base(path)
			if ok, err := filepath.Match("*.mp3", base); !ok || err != nil {
				if err != nil {
					return err
				}
				return nil
			}

			dir := filepath.Dir(path)
			if currdir != dir {
				if currdir != "" && len(files) > 0 {
					listtags(currdir, files, maxlen)
					files = []string{}
					maxlen = 0
				}
				currdir = dir
			}

			if info.IsDir() {
				return nil
			}

			files = append(files, base)
			if len(base) > maxlen {
				maxlen = len(base)
			}

			return nil
		})
		listtags(currdir, files, maxlen)
	}
}

func listtags(path string, files []string, maxlen int) {
	if path == "" || len(files) == 0 {
		return
	}

	filenamespec := fmt.Sprint("   %-", maxlen, "s")

	println()
	println(path)
	for _, filename := range files {
		fmt.Printf(filenamespec, filename)
		mp3, err := mp3.FromFile(filepath.Join(path, filename))
		if err != nil {
			fmt.Printf("\nERROR: %s\n", err)
			return
		}

		if mp3.Id3v1 != nil {
			fmt.Printf("     %s", mp3.Id3v1.Version)
		}

		for _, v2 := range mp3.Id3v2 {
			fmt.Printf("    %s", v2.Version)
		}

		println()
	}
}
