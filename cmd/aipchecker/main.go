package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/spin"
	"github.com/wozuo/aipchecker/checker"
	"github.com/wozuo/aipchecker/zipper"
)

var usageMsg = `usage: aipchecker [flags] [path/to/Android/projects/folder]
Flags:
	The -unzip flag indicates if projects need to be unzipped first
	
Examples:
	aipchecker -unzip "~/documents/Android Projects"
	aipchecker "~/documents/Android Projects"
`

func usage() {
	fmt.Fprint(os.Stderr, usageMsg)
	os.Exit(2)
}

func main() {
	zf := flag.Bool("unzip", false, "unzip projects first")
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) == 3 && *zf {
		unzipProjects(os.Args[2])
		checkPermissions(os.Args[2])
		fmt.Println("All done!")
		return
	} else if len(os.Args) == 2 {
		checkPermissions(os.Args[1])
		fmt.Println("All done!")
		return
	}

	fmt.Println("Invalid arguments.")
}

func unzipProjects(path string) {
	spinner := spin.New("%s Unzipping...")
	spinner.Start()

	err := zipper.UnzipAll(path)

	if err != nil {
		fmt.Println("Error unzipping: ", err)
		os.Exit(3)
	}

	spinner.Stop()
}

func checkPermissions(path string) {
	spinner := spin.New("%s Checking permissions...")
	spinner.Start()

	checker.CheckProjects(path)
	spinner.Stop()
}
