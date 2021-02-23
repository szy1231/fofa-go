// package main
package fofa_cli

import (
	"flag"
	"os"
)

func init() {
	flag.Usage = usage
	if len(os.Args) <= 1 {
		usage()
		os.Exit(0)
	}
	parse()
}

func main() {

	if *help {
		usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "version":
		cliVersion()
	case "init":
		userInit()
		os.Exit(0)
	case "info":
		setClt()
		getInfo()
	case "search":
		setClt()
		search()
	default:
		usage()
		os.Exit(0)
	}

}
