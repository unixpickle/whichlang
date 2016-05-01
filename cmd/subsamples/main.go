package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <num lines> <dir>\n", os.Args[0])
		os.Exit(1)
	}
	numLines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	sampleDir := os.Args[2]
	if err := subsample(numLines, sampleDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func subsample(numLines int, dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	listing, err := dir.Readdir(-1)
	dir.Close()
	if err != nil {
		return err
	}
	for _, langDir := range listing {
		if !langDir.IsDir() {
			continue
		}
		p := filepath.Join(dirPath, langDir.Name())
		if err := subsampleLang(numLines, p); err != nil {
			return err
		}
	}
	return nil
}

func subsampleLang(numLines int, langDirPath string) error {
	dir, err := os.Open(langDirPath)
	if err != nil {
		return err
	}
	listing, err := dir.Readdir(-1)
	dir.Close()
	if err != nil {
		return err
	}
	for _, fileInfo := range listing {
		if strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		filePath := filepath.Join(langDirPath, fileInfo.Name())
		if err := subsampleFile(numLines, filePath); err != nil {
			return err
		}
	}
	return nil
}

func subsampleFile(numLines int, filePath string) error {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) < numLines {
		fmt.Println("Skipping file", filePath)
		return nil
	}
	startIndex := (len(lines) - numLines) / 2
	splitLines := lines[startIndex : startIndex+numLines]
	newPath := newPath(numLines, filePath)
	newData := strings.Join(splitLines, "\n")
	return ioutil.WriteFile(newPath, []byte(newData), 0755)
}

func newPath(numLines int, filePath string) string {
	sampleTag := "_subsample_" + strconv.Itoa(numLines)
	ext := filepath.Ext(filePath)
	if ext == "" {
		return filePath + sampleTag
	} else {
		return filePath[:len(filePath)-len(ext)] + sampleTag + ext
	}
}
