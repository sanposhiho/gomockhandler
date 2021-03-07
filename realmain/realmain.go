package realmain

import "github.com/sanposhiho/gomockhandler/model"

type ChunkRepo interface {
	Put(m *model.Chunk) error
	Get() (*model.Chunk, error)
}

type MockgenRunner interface {
	Run() error
}

type Runner struct {
	ChunkRepo     ChunkRepo
	MockgenRunner MockgenRunner

	Args Args
}

type Args struct {
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
