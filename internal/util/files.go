package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Tree(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get dir: %w", err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			childPaths, err := Tree(filepath.Join(dir, file.Name()))
			if err != nil {
				return []string{}, err
			}
			paths = append(paths, childPaths...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}

func ReadALine(filename string) (string, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return "", fmt.Errorf("there is not %s: %w", filename, err)
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "", fmt.Errorf("error to open the file: %w", err)
	}

	var line string
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line = fileScanner.Text()
		break
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("error to scan the file: %s", err)
	}

	return line, nil
}
