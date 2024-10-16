package mockgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sanposhiho/gomockhandler/internal/util"
)

type Runner interface {
	Run() error

	SetSource(new string)
	GetSource() string
	SetDestination(new string)
	GetDestination() string
}

func Checksum(r Runner) (string, error) {
	d := r.GetDestination()
	tmpFilePath := util.TmpFilePath(d)
	defer os.Remove(tmpFilePath)

	// use tmpfile to test generating mock
	r.SetDestination(tmpFilePath)
	defer r.SetDestination(d)

	if err := r.Run(); err != nil {
		return "", fmt.Errorf("failed to run mockgen: %v \nPlease run `%s` and check if mockgen works correctly with your options", err, r)
	}

	tmpFile, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed read file. filename: %s, err: %w", tmpFilePath, err)
	}

	// See https://github.com/sanposhiho/gomockhandler/issues/88
	adjustedFile := replaceTmpPathWithOriginal(tmpFilePath, tmpFile)

	checksum, err := util.CalculateCheckSum([]byte(adjustedFile))
	if err != nil {
		return "", fmt.Errorf("calculate checksum of the mock: %v", err)
	}

	return checksum, nil
}

func replaceTmpPathWithOriginal(tmpFilePath string, tmpFile []byte) string {
	d, f := filepath.Split(tmpFilePath)
	originFilePath := filepath.Join(d, strings.Replace(f, "tmp_", "", 1))
	return strings.Replace(string(tmpFile), tmpFilePath, originFilePath, 1)
}

func SourceChecksum(r Runner) (string, error) {
	file, err := os.ReadFile(r.GetSource())
	if err != nil {
		return "", fmt.Errorf("failed read file. filename: %s, err: %w", r.GetSource(), err)
	}
	checksum, err := util.CalculateCheckSum(file)
	if err != nil {
		return "", fmt.Errorf("calculate checksum of the mock source: %v", err)
	}

	return checksum, nil
}
