package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func CalculateCheckSum(filePath string) (string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed read file. filename: %s, err: %w", filePath, err)
	}

	originFilePath := OriginFilePathFromTmpFilePath(filePath)
	re := regexp.MustCompile("(//.+mockgen.+-destination=)(" + filePath + ")")
	trimmedFile := re.ReplaceAllString(string(file), "${1}"+originFilePath)

	hash := md5.Sum([]byte(trimmedFile))
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

func OriginFilePathFromTmpFilePath(tmpFilePath string) string {
	return strings.Replace(tmpFilePath, "tmp_", "", 1)
}
