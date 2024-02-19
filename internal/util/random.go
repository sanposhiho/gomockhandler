package util

import (
	"crypto/md5"
	"encoding/base64"
	"path/filepath"
	"strings"
)

func CalculateCheckSum(file []byte) (string, error) {
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
