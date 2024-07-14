package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	command := os.Args
	if len(command) > 1 {
		path := command[1]
		err := walkDirectory(path)
		if err != nil {
			println(err)
		}
	}
}

func scanFile(file *os.File) ([]string, error) {
	scanner := bufio.NewScanner(file)

	line := 1
	found := false
	var foundLines []string

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "TODO:") {
			found = true
			foundLines = append(foundLines, strconv.Itoa(line))
		} else if found && len(scanner.Text()) == 0 {
			found = false
		}
		if found {
			foundLines = append(foundLines, strings.Trim(scanner.Text(), " "))
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		return foundLines, err
	}

	return foundLines, nil

}

func readFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	foundLine, err := scanFile(file)
	if err != nil {
		return err
	}

	if len(foundLine) > 0 {

		fmt.Println("=>", path)
		fmt.Println("-----------------------------------------------------------------------")

		for i := 0; i < len(foundLine); i++ {

			fmt.Println(" ", foundLine[i])
		}

		fmt.Println("-----------------------------------------------------------------------")
	}
	return nil

}

func walkDirectory(path string) error {
	err := filepath.WalkDir(path, func(path string, res fs.DirEntry, err error) error {
		// fmt.Println(res.Name(), res.IsDir())
		if !res.IsDir() {
			err := readFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil

}
