package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type FileContent struct {
	Path    string
	Content []byte
}

var filecountParse int

func CheckImport(metadatafolder string, datafolder string) error {
	files, err := walkFiles(metadatafolder, "artifactory-file.xml")
	fmt.Printf("\n")
	if err != nil {
		return err
	}

	inprogress := false

	fmt.Printf("Found metadata files: %d\n", len(files))

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	for _, file := range files {
		var artfile ArtifactoryFile
		if err := xml.Unmarshal(file.Content, &artfile); err != nil {
			if inprogress {
				fmt.Printf("\n")
				inprogress = false
			}
			fmt.Printf("failed to parse xml file: '%s': %v", file.Path, err)
			continue
		}

		var sha1sum string
		for _, checksum := range artfile.AdditionalInfo.ChecksumsInfo.Checksums.Checksum {
			if checksum.Type == "sha1" {
				if len(checksum.Actual) != 40 {
					if inprogress {
						fmt.Printf("\n")
						inprogress = false
					}
					fmt.Printf("'%s' (%d): warning: sha1 is not 40 chars: '%s'\n", file.Path, artfile.Size, checksum.Actual)
				}
				sha1sum = checksum.Actual
				break
			}
		}

		datapath := filepath.Join(datafolder, sha1sum[0:2], sha1sum)

		info, err := os.Stat(datapath)
		if err != nil {
			if os.IsNotExist(err) {
				if inprogress {
					fmt.Printf("\n")
					inprogress = false
				}
				fmt.Printf("'%s' (%d): warning: data file not found: '%s'\n", file.Path, artfile.Size, datapath)
				continue
			}
		}
		if info.IsDir() {
			if inprogress {
				fmt.Printf("\n")
				inprogress = false
			}
			fmt.Printf("warning: '%s' is a directory, not a file.\n", datapath)
			continue
		}

		filecountParse++
		if filecountParse%10000 == 0 {
			fmt.Printf(".")
			inprogress = true
		}
	}
	if inprogress {
		fmt.Printf("\n")
		inprogress = false
	}

	return nil
}

var filecountWalk int

func walkFiles(root, pattern string) ([]FileContent, error) {
	var files []FileContent

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		filecountWalk++
		if filecountWalk%10000 == 0 {
			fmt.Printf(".")
		}

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		match, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if match {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, FileContent{
				Path:    relPath,
				Content: content,
			})
		}

		return nil
	})

	return files, err
}
