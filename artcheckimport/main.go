package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	reverseCheck := flag.Bool("f", false, "Perform a reverse check.")

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 || args[0] == "" || args[1] == "" {
		usage()
		os.Exit(1)
	}

	fmt.Printf("reverseCheck: %t\n", *reverseCheck)

	metadatafolder := args[0]
	datafolder := args[1]

	if !folderExists(metadatafolder) {
		fmt.Printf("Error: metadata folder '%s' does not exist or is not a directory.\n", metadatafolder)
		os.Exit(1)
	}
	if !folderExists(datafolder) {
		fmt.Printf("Error: data folder '%s' does not exist or is not a directory.\n", datafolder)
		os.Exit(1)
	}

	err := CheckImport(metadatafolder, datafolder)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

func usage() {
	fmt.Println("ARTCHECKIMPORT - Artifactory import checker tool")
	fmt.Println()
	fmt.Println("This tool is used for verifying that all artifacts exists, before system import.")
	fmt.Println()
	fmt.Println("Usage: artcheckimport [-r] <metadatafolder> <datafolder>")
	fmt.Println()
	flag.PrintDefaults()
}
