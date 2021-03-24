package command

import (
	"github.com/sanposhiho/gomockhandler/internal/mockgen"
	"github.com/sanposhiho/gomockhandler/internal/model"
)

type ConfigRepo interface {
	Put(m *model.Config, path string) error
	Get(path string) (*model.Config, error)
}

type Runner struct {
	ConfigRepo    ConfigRepo
	MockgenRunner mockgen.Runner

	Args Args
}

type Args struct {
	ConfigPath string

	Source          string
	Destination     string
	MockNames       string
	PackageOut      string
	SelfPackage     string
	WritePkgComment bool
	CopyrightFile   string
	Imports         string
	AuxFiles        string
	ExecOnly        string
	BuildFlags      string
	ProgOnly        bool
	DebugParser     bool
}
