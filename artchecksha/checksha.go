package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var totalFiles int
var wrongFiles int
var totalBytes int64
var startTime time.Time

func CheckSHA1(datafolder string, start string, end string) error {
	printLeadingNewline := false
	err := filepath.Walk(datafolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		parentDir := filepath.Base(filepath.Dir(path))

		if isValidHexByte(start) && isValidHexByte(end) {
			if !isValidHexByte(parentDir) {
				return nil
			}
			if parentDir < start || parentDir > end {
				return nil
			}
		}

		totalFiles++
		totalBytes += info.Size()

		if totalFiles%10 == 0 {
			duration := time.Since(startTime)

			mb := float64(totalBytes) / (1024 * 1024)
			seconds := duration.Seconds()
			throughput := mb / seconds

			fmt.Printf("Total: %d, Wrong: %d, Correct: %d, Duration: %s, Bytes: %d, Throughput: %.2f MB/s, Folder: '%s'.          \r",
				totalFiles, wrongFiles, totalFiles-wrongFiles, duration, totalBytes, throughput, parentDir)
			printLeadingNewline = true
		}

		hash, err := computeSHA1(path)
		if err != nil {
			if printLeadingNewline {
				fmt.Println()
				printLeadingNewline = false
			}
			fmt.Printf("Error hashing %s: %v\n", path, err)
			return nil
		}

		if info.Name() != hash {
			if printLeadingNewline {
				fmt.Println()
				printLeadingNewline = false
			}
			fmt.Printf("Error: expected (filename/sha1): '%s', actual: '%s'\n", path, hash)
			wrongFiles++
		}

		return nil
	})

	if printLeadingNewline {
		fmt.Println()
	}

	return err
}

func computeSHA1(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func PrintStats() {
	duration := time.Since(startTime)

	mb := float64(totalBytes) / (1024 * 1024)
	seconds := duration.Seconds()
	throughput := mb / seconds

	fmt.Println("---------- Statistics ----------")
	fmt.Printf("Total files checked: %d\n", totalFiles)
	fmt.Printf("Files with wrong hash: %d\n", wrongFiles)
	fmt.Printf("Files with correct hash: %d\n", totalFiles-wrongFiles)
	fmt.Printf("Total duration: %s\n", duration)
	fmt.Printf("Total bytes: %d\n", totalBytes)
	fmt.Printf("Throughput: %2f MB/s\n", throughput)
	fmt.Println("--------------------------------")
}
