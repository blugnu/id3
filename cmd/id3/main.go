package main

import (
	"flag"
	"os"
)

var listcmd bool
var showcmd bool

func main() {
	var showusage bool

	flag.BoolVar(&listcmd, "list", false, "list tags identified in files")
	flag.BoolVar(&showcmd, "show", false, "show tags identified in files")
	flag.BoolVar(&showusage, "usage", false, "display this usage information")
	flag.Parse()

	if len(os.Args) == 1 || showusage {
		flag.Usage()
	}

	if showcmd {
		show()
	} else if listcmd {
		list()
	}

}
