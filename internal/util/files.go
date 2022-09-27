package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Tree(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get dir: %w", err)
	}

	var paths []string
	for _, file := range files {
		if strings.Contains(file.Name(), "vendor") {
			// vendor dir should be ignored.
			continue
		}

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

	r := bufio.NewReader(file)

	data, cont, err := r.ReadLine()
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error to read the file: %w", err)
	}

	buf := bytes.NewBuffer(data)
	for cont {
		data, cont, err = r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("error to read the file: %w", err)
		}

		if _, err = buf.Write(data); err != nil {
			return "", fmt.Errorf("error to append data to buffer: %w", err)
		}
	}

	return buf.String(), nil
}
