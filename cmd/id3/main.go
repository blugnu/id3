package main

import (
	"flag"
	"os"
	"strings"
)

type _cmd struct {
	clean bool
	list  bool
	show  bool
}

var cmd = _cmd{}

func (_cmd) String() string {
	cmds := []string{}
	if cmd.clean {
		cmds = append(cmds, "--clean")
	}
	if cmd.list {
		cmds = append(cmds, "--list")
	}
	if cmd.show {
		cmds = append(cmds, "--show")
	}
	return strings.Join(cmds, ",")
}

func (_cmd) run() {
	dolist := cmd.list || (cmd.String() == (_cmd{}).String() && len(flag.Args()) > 0)

	if dolist {
		list()
	}
	if cmd.show {
		show()
	}
	if cmd.list {
		list()
	}
}

func main() {
	var showusage bool

	flag.BoolVar(&cmd.clean, "clean", false, "cleanup tags (remove unwanted tags and values) and set values from the filename (where possible)")
	flag.BoolVar(&cmd.list, "list", false, "list tags identified in files")
	flag.BoolVar(&cmd.show, "show", false, "show tags identified in files")
	flag.BoolVar(&showusage, "usage", false, "display this usage information")
	flag.Parse()

	if len(os.Args) == 1 || showusage {
		flag.Usage()
	}

	cmd.run()
}
