package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type InvalidArgsError struct{}
type InvalidSourceArg struct{}
type InvalidGlobArg struct{}

func (e InvalidArgsError) Error() string {
	return "Invalid amount of arguments"
}

func (e InvalidSourceArg) Error() string {
	return "Invalid source arg"
}

func (e InvalidGlobArg) Error() string {
	return "Invalid glob arg"
}

func searchAndReplace() error {
	if len(os.Args) != 5 {
		return InvalidArgsError{}
	}

	sourceDir := os.Args[1]
	globPattern := os.Args[2]
	stringToReplaceFrom := os.Args[3]
	stringToReplaceTo := os.Args[4]

	if sourceDir == "" {
		return InvalidSourceArg{}
	}

	if globPattern == "" {
		return InvalidGlobArg{}
	}

	err := filepath.WalkDir(sourceDir, func(path string, dir os.DirEntry, err error) error {
		if dir.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		isMatch, err := filepath.Match(globPattern, filepath.Base(path))

		if err != nil {
			return err
		}

		if !isMatch {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		ext := filepath.Ext(path)

		modifiedFileName := fmt.Sprintf("%s-modified%s", file.Name(), ext)
		modifiedFile, err := os.Create(modifiedFileName)
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if errors.Is(err, io.EOF) {
				lineToWrite := line
				if strings.Contains(line, stringToReplaceFrom) {
					lineToWrite = strings.Replace(line, stringToReplaceFrom, stringToReplaceTo, 1)
				}

				if _, err := modifiedFile.WriteString(lineToWrite); err != nil {
					return err
				}

				break
			}

			if err != nil {
				removeErr := os.RemoveAll(modifiedFileName)
				return errors.New(removeErr.Error() + "\n" + err.Error())
			}

			lineToWrite := line
			if strings.Contains(line, stringToReplaceFrom) {
				lineToWrite = strings.ReplaceAll(line, stringToReplaceFrom, stringToReplaceTo)
			}

			if _, err = modifiedFile.WriteString(lineToWrite); err != nil {
				return err
			}
		}

		if err := os.Remove(file.Name()); err != nil {
			return err
		}

		if err := os.Rename(modifiedFile.Name(), file.Name()); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := searchAndReplace(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
