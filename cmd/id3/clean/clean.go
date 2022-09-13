package clean

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v2"
	"github.com/blugnu/tags/mp3"
)

func parseInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 32)
	return int(v)
}

var rxiTunes = regexp.MustCompile(`^([0-9]*)\. (.*)\.mp3$`)
var rxSingleDisc = regexp.MustCompile(`^(.+) - ([0-9]+) - (.+)\.mp3$`)
var rxMultiDisc = regexp.MustCompile(`^(.+) - ([0-9]+)-([0-9]+) - (.+)\.mp3$`)

type fileinfo struct {
	filename string
	album    string
	title    string
	discno   int
	trackno  int
}

func Run() {
	for _, arg := range flag.Args() {
		currdir := ""
		files := []fileinfo{}
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
					cleanfiles(currdir, files)
					files = []fileinfo{}
				}
				currdir = dir
			}

			if info.IsDir() {
				return nil
			}

			f := fileinfo{
				filename: base,
			}

			el := rxiTunes.FindStringSubmatch(base)
			if el == nil {
				el = rxSingleDisc.FindStringSubmatch(base)
			}
			if el == nil {
				el = rxMultiDisc.FindStringSubmatch(base)
			}
			switch len(el) {
			case 3:
				f.album = filepath.Base(currdir)
				f.discno = 1
				f.trackno = parseInt(el[1])
				f.title = el[2]
			case 4:
				f.discno = 1
				f.album = el[1]
				f.trackno = parseInt(el[2])
				f.title = el[3]
			case 5:
				f.album = el[1]
				f.discno = parseInt(el[2])
				f.trackno = parseInt(el[3])
				f.title = el[4]
			}
			files = append(files, f)

			return nil
		})
		cleanfiles(currdir, files)
	}
}

func cleanfiles(path string, files []fileinfo) {
	maxdisc := 1
	maxtrack := 1
	for _, f := range files {
		if f.discno > maxdisc {
			maxdisc = f.discno
		}
		if f.trackno > maxtrack {
			maxtrack = f.trackno
		}
	}

	outdir := filepath.Join(filepath.Dir(path), "clean")
	println("output to ", outdir)

	for _, f := range files {
		mp3, err := mp3.FromFile(filepath.Join(path, f.filename))
		if err != nil {
			println(err)
			continue
		}

		v1 := mp3.Id3v1
		v22 := mp3.GetTag(id3.Id3v22)
		v23 := mp3.GetTag(id3.Id3v23)
		v24 := mp3.GetTag(id3.Id3v24)

		newtag := v23 == nil
		if newtag {
			v23 = mp3.CreateTag(id3.Id3v23)
		}

		var keep = id3.FrameKeySet{id3.TALB, id3.TCON, id3.TCOM, id3.TIT2, id3.TPE1, id3.TPE2, id3.TPOS, id3.TRCK, id3.TCMP, id3.TYER, id3.APIC}
		copy := keep.Remove(id3.APIC)

		retained := []*id3v2.Frame{}
		for _, frame := range v23.Frames {
			if keep.Contains(frame.Key) {
				retained = append(retained, frame)
			}
		}
		v23.Frames = retained

		for _, key := range copy {
			value := v23.Get(key)
			if v1 != nil {
				value1 := v1.Get(key)
				if len(value1) > len(value) {
					value = value1
				}
			}
			if v22 != nil {
				value22 := v22.Get(key)
				if len(value22) > len(value) {
					value = value22
				}
			}
			if v24 != nil {
				value24 := v24.Get(key)
				if len(value24) > len(value) {
					value = value24
				}
			}
			v23.Set(key, value)
		}

		mp3.Id3v1 = nil
		mp3.Id3v2 = []*id3v2.Tag{v23}

		v23.Set(id3.TALB, f.album)
		v23.Set(id3.TIT2, f.title)
		v23.SetInt(id3.NumDiscs, maxdisc)
		v23.SetInt(id3.NumTracks, maxtrack)

		for _, v2 := range mp3.Id3v2 {
			println()
			tagver := fmt.Sprintf("- %s", v2.Version.String())
			for _, frame := range v2.Frames {
				value, ok := frame.Data.(string)
				if !ok {
					var stringer fmt.Stringer
					stringer, ok = frame.Data.(fmt.Stringer)
					if ok {
						value = stringer.String()
					}
				}
				if !ok {
					fmt.Printf("    %-8s  %v  (%d bytes, %T)\n", tagver, frame.ID, frame.Size, frame.Data)
				} else {
					fmt.Printf("    %-8s  %v  %s\n", tagver, frame.ID, value)
				}
				tagver = ""
			}
		}
	}
}
