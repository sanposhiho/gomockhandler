package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func CalculateCheckSum(filePath string) (string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed read file. filename: %s, err: %w", filePath, err)
	}

	hash := md5.Sum(file)
	strhash := base64.StdEncoding.EncodeToString(hash[:])
	return strhash, nil
}

func PathInProject(projectRoot, path string) string {
	return filepath.Clean(strings.Replace(path, projectRoot, "", 1))
}

func TmpFilePath(original string) string {
	d, f := filepath.Split(original)
	return d + "tmp_" + f
}
