package mockgen

import (
	"fmt"
	"os"

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
	tmpFile := util.TmpFilePath(d)
	defer os.Remove(tmpFile)

	// use tmpfile to test generating mock
	r.SetDestination(tmpFile)
	defer r.SetDestination(d)

	if err := r.Run(); err != nil {
		return "", fmt.Errorf("failed to run mockgen: %v \nPlease run `%s` and check if mockgen works correctly with your options", err, r)
	}

	checksum, err := util.CalculateCheckSum(tmpFile)
	if err != nil {
		return "", fmt.Errorf("calculate checksum of the mock: %v", err)
	}

	return checksum, nil
}

func SourceChecksum(r Runner) (string, error) {
	checksum, err := util.CalculateCheckSum(r.GetSource())
	if err != nil {
		return "", fmt.Errorf("calculate checksum of the mock source: %v", err)
	}

	return checksum, nil
}
