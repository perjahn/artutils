package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	datafolder := os.Args[1]

	if !folderExists(datafolder) {
		fmt.Printf("Error: folder '%s' does not exist or is not a directory.\n", datafolder)
		os.Exit(1)
	}

	startTime = time.Now()
	start := ""
	end := ""
	if len(os.Args) == 4 {
		start = os.Args[2]
		end = os.Args[3]
		if !isValidHexByte(start) || !isValidHexByte(end) {
			fmt.Printf("Error: Start and End must be valid hex digits (i.e., '00' to 'ff').\n")
			os.Exit(1)
		}
	}
	err := CheckSHA1(datafolder, start, end)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	PrintStats()
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

func isValidHexByte(s string) bool {
	if len(s) != 2 {
		return false
	}
	if ((s[0] < '0' || s[0] > '9') && (s[0] < 'a' || s[0] > 'f')) ||
		((s[1] < '0' || s[1] > '9') && (s[1] < 'a' || s[1] > 'f')) {
		return false
	}
	return true
}

func usage() {
	fmt.Println("ARTCHECKSHA - Artifactory hash checker tool")
	fmt.Println()
	fmt.Println("This tool is used for verifying that all artifacts have correct sha1 hash sum.")
	fmt.Println()
	fmt.Println("Usage: artchecksha <folder> [start] [end]")
	fmt.Println()
	flag.PrintDefaults()
}
