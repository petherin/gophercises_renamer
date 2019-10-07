package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type fileData struct {
	fileName string
	number   string
	ext      string
}

const root = "sample"

func main() {
	filesToRename := map[string][]fileData{}
	group := ""

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		splitOne := strings.Split(path, "_")
		if len(splitOne) < 2 {
			return nil
		}

		if group != splitOne[0] {
			filesToRename[splitOne[0]] = []fileData{}
			group = splitOne[0]
		}

		// get number and file extension
		numberExt := strings.Split(splitOne[1], ".")
		if len(numberExt) == 0 {
			return nil
		}

		fileNumber, err := strconv.Atoi(numberExt[0])
		if err != nil {
			fmt.Printf("%s unexpected format", path)
			return nil
		}

		currentFile := fileData{fileName: path, number: strconv.Itoa(fileNumber), ext: numberExt[1]}

		// append to array in group's map entry
		filesToRename[group] = append(filesToRename[group], currentFile)

		return nil
	})

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// Loop through map of files to rename. In each map is an array of files with the same base name.
	// Loop through the arrays backwards so we know the max number of each file.
	for pathKey, fileGroup := range filesToRename {
		maxNum := ""
		for i := len(fileGroup) - 1; i >= 0; i-- {
			if i == len(fileGroup)-1 {
				maxNum = fileGroup[i].number
			}

			newName := fmt.Sprintf("%s (%s of %s).txt", pathKey, fileGroup[i].number, maxNum)
			err := os.Rename(fileGroup[i].fileName, newName)
			if err != nil {
				fmt.Printf("error occured: %s\n", err)
			}
		}
	}
}
