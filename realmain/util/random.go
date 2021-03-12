package util

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func MockCheckSum(filePath string) ([16]byte, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed read file. filename: %s, err: %w", filePath)
	}

	hash := md5.Sum(file)
	return hash, nil
}

func PathInProject(projectRoot, path string) string {
	return filepath.Clean(strings.Replace(path, projectRoot, "", 1))
}
