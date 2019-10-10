package main

import (
	"bufio"
	"flag"
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
	safe := flag.Bool("safe", true, "when true, no files will be renamed, just a log of what would have happened")
	flag.Parse()

	safeMsg :=""
	if *safe {
		fmt.Print("In safe mode, no files will be renamed.\n")
		safeMsg = " (safe mode)"
	} else {
		fmt.Print("Not in safe mode, files will be renamed. Proceed (y/n)?\n")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		switch char {
		case 'y', 'Y':
			fmt.Println("Starting rename...\n")
		default:
			return
		}
	}

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
			fmt.Printf("Renaming %s to %s%s\n", fileGroup[i].fileName, newName, safeMsg)

			if *safe {
				continue
			}

			err := os.Rename(fileGroup[i].fileName, newName)
			if err != nil {
				fmt.Printf("error occured: %s\n", err)
			}
		}
	}
}
