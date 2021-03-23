package command

import (
	"github.com/sanposhiho/gomockhandler/internal/mockgen"
	"github.com/sanposhiho/gomockhandler/internal/model"
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
	ConfigPath  string
	Concurrency int

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
