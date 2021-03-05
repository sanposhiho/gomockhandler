package mockgen

import (
	"os/exec"
	"strconv"
)

type mode int

const (
	SourceMode mode = iota
	ReflectMode
)

type Runner struct {
	Source          string
	Destination     string
	Package         string
	Imports         string
	AuxFiles        string
	BuildFlags      string
	MockNames       string
	SelfPackage     string
	CopyrightFile   string
	execOnly        string
	progOnly        bool
	writePkgComment bool
	debugParser     bool
}

func NewRunner(source, dest, pkg, imp, af, bf, mn, spkg, cf, eo string, po, wpc, dp bool) *Runner {
	return &Runner{
		Source:          source,
		Destination:     dest,
		Package:         pkg,
		Imports:         imp,
		AuxFiles:        af,
		BuildFlags:      bf,
		MockNames:       mn,
		SelfPackage:     spkg,
		CopyrightFile:   cf,
		execOnly:        eo,
		progOnly:        po,
		writePkgComment: wpc,
		debugParser:     dp,
	}
}

func (r *Runner) Run() error {
	return exec.Command("mockgen", r.options()...).Run()
}

func (r *Runner) options() []string {
	var opts []string
	if r.Source != "" {
		opts = append(opts, []string{
			"-source", r.Source,
		}...)
	}
	if r.Destination != "" {
		opts = append(opts, []string{
			"-destination", r.Destination,
		}...)
	}
	if r.Package != "" {
		opts = append(opts, []string{
			"-package", r.Package,
		}...)
	}
	if r.Imports != "" {
		opts = append(opts, []string{
			"-imports", r.Imports,
		}...)
	}
	if r.AuxFiles != "" {
		opts = append(opts, []string{
			"-aux_files", r.AuxFiles,
		}...)
	}
	if r.BuildFlags != "" {
		opts = append(opts, []string{
			"-build_flags", r.BuildFlags,
		}...)
	}
	if r.MockNames != "" {
		opts = append(opts, []string{
			"-mock_names", r.MockNames,
		}...)
	}
	if r.SelfPackage != "" {
		opts = append(opts, []string{
			"-self_package", r.SelfPackage,
		}...)
	}
	if r.CopyrightFile != "" {
		opts = append(opts, []string{
			"-copyright_file", r.CopyrightFile,
		}...)
	}
	if r.execOnly != "" {
		opts = append(opts, []string{
			"-exec_only", r.execOnly,
		}...)
	}

	opts = append(opts, []string{
		"-prog_only", strconv.FormatBool(r.progOnly),
	}...)
	opts = append(opts, []string{
		"-write_package_comment", strconv.FormatBool(r.writePkgComment),
	}...)
	opts = append(opts, []string{
		"-debug_parser", strconv.FormatBool(r.writePkgComment),
	}...)

	return opts
}
