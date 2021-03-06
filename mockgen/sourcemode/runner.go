package sourcemode

import (
	"os/exec"
	"strconv"
)

type Runner struct {
	Source          string
	Destination     string
	Package         string
	Imports         string
	AuxFiles        string
	MockNames       string
	SelfPackage     string
	CopyrightFile   string
	WritePkgComment bool
	DebugParser     bool
}

func NewRunner(source, dest, pkg, imp, af, mn, spkg, cf string, wpc, dp bool) *Runner {
	return &Runner{
		Source:          source,
		Destination:     dest,
		Package:         pkg,
		Imports:         imp,
		AuxFiles:        af,
		MockNames:       mn,
		SelfPackage:     spkg,
		CopyrightFile:   cf,
		WritePkgComment: wpc,
		DebugParser:     dp,
	}
}

func (r *Runner) Run() error {
	return exec.Command("mockgen", r.options()...).Run()
}

func (r *Runner) options() []string {
	var opts []string
	if r.Source != "" {
		opts = append(opts, "-source="+r.Source)
	}
	if r.Destination != "" {
		opts = append(opts, "-destination="+r.Destination)
	}
	if r.Package != "" {
		opts = append(opts, "-package="+r.Package)
	}
	if r.Imports != "" {
		opts = append(opts, "-imports="+r.Imports)
	}
	if r.AuxFiles != "" {
		opts = append(opts, "-aux_files="+r.AuxFiles)
	}
	if r.MockNames != "" {
		opts = append(opts, "-mock_names="+r.MockNames)
	}
	if r.SelfPackage != "" {
		opts = append(opts, "-self_package="+r.SelfPackage)
	}
	if r.CopyrightFile != "" {
		opts = append(opts, "-copyright_file="+r.CopyrightFile)
	}
	opts = append(opts, "-write_package_comment="+strconv.FormatBool(r.WritePkgComment))
	opts = append(opts, "-debug_parser="+strconv.FormatBool(r.DebugParser))

	return opts
}
