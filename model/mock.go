package model

import (
	"github.com/sanposhiho/gomockhandler/mockgen"
	"github.com/sanposhiho/gomockhandler/mockgen/reflectmode"
	"github.com/sanposhiho/gomockhandler/mockgen/sourcemode"
)

type mode string

const (
	Unknown     mode = "UNKNOWN"
	ReflectMode mode = "REFLECT_MODE"
	SourceMode  mode = "SOURCE_MODE"
)

type Mock struct {
	Destination       string              `json:"destination"`
	CheckSum          [16]byte            `json:"checksum"`
	Mode              mode                `json:"mode"`
	ReflectModeRunner *reflectmode.Runner `json:"reflect_mode_runner,omitempty"`
	SourceModeRunner  *sourcemode.Runner  `json:"source_mode_runner,omitempty"`
}

func NewMock(destination string, checksum [16]byte, genrunner mockgen.Runner) Mock {
	rrunner, srunner := convertMockgenRunner(genrunner)
	mode := Unknown
	if rrunner != nil {
		mode = ReflectMode
	} else if srunner != nil {
		mode = SourceMode
	}
	return Mock{
		Destination:       destination,
		CheckSum:          checksum,
		Mode:              mode,
		ReflectModeRunner: rrunner,
		SourceModeRunner:  srunner,
	}
}

func convertMockgenRunner(r mockgen.Runner) (*reflectmode.Runner, *sourcemode.Runner) {
	switch runner := r.(type) {
	case *reflectmode.Runner:
		return runner, nil
	case *sourcemode.Runner:
		return nil, runner
	}
	return nil, nil
}
