package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/blugnu/tags/mp3"
)

func show() {
	for _, arg := range flag.Args() {
		currdir := ""
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
				println()
				println(dir)
				currdir = dir
			}

			showtags(path)

			return nil
		})
	}
}

func showtags(path string) {
	println()
	println("   ", filepath.Dir(path))
	println("   ", filepath.Base(path))

	mp3, err := mp3.FromFile(path)
	if err != nil {
		return
	}

	v1 := mp3.Id3v1
	if v1 != nil {
		println()
		fmt.Printf("    - %-6s  Album    %s\n", v1.Version, v1.Album)
		fmt.Printf("              Artist   %s\n", v1.Artist)
		fmt.Printf("              Title    %s\n", v1.Title)
		fmt.Printf("              Year     %d\n", v1.Year)
		fmt.Printf("              Genre    %s\n", v1.Genre)
		fmt.Printf("              Comment  %s\n", v1.Comment)
	}

	for _, v2 := range mp3.Id3v2 {
		println()
		tagver := fmt.Sprintf("- %s", v2.Version.String())
		for _, frame := range v2.Frames {
			if frame.Text != nil {
				fmt.Printf("    %-8s  %s  %s\n", tagver, frame.ID, *frame.Text)
			} else {
				fmt.Printf("    %-8s  %s  %d bytes\n", tagver, frame.ID, frame.Size)
			}
			tagver = ""
		}
	}

	println()
}
