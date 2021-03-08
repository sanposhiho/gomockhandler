package realmain

import (
	"github.com/sanposhiho/gomockhandler/mockgen"
	"github.com/sanposhiho/gomockhandler/model"
)

type ChunkRepo interface {
	Put(m *model.Config, path string) error
	Get(path string) (*model.Config, error)
}

type Runner struct {
	ChunkRepo     ChunkRepo
	MockgenRunner mockgen.Runner

	Args Args
}

type Args struct {
	ProjectRoot string

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
