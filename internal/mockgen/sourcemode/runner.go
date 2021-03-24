package sourcemode

import (
	"os/exec"
	"strconv"
)

type Runner struct {
	Source          string `json:"source,omitempty"`
	Destination     string `json:"destination,omitempty"`
	Package         string `json:"package,omitempty"`
	Imports         string `json:"imports,omitempty"`
	AuxFiles        string `json:"aux_files,omitempty"`
	MockNames       string `json:"mock_names,omitempty"`
	SelfPackage     string `json:"self_package,omitempty"`
	CopyrightFile   string `json:"copyright_file,omitempty"`
	WritePkgComment *bool  `json:"write_pkg_comment,omitempty"`
	DebugParser     *bool  `json:"debug_parser,omitempty"`
}

func NewRunner(source, dest, pkg, imp, af, mn, spkg, cf string, wpc, dp bool) *Runner {
	var wpcp *bool
	if wpc != true {
		// The default value of wpc is true
		wpcp = &wpc
	}

	var dpp *bool
	if dp != false {
		// The default value of dp is false
		dpp = &dp
	}
	return &Runner{
		Source:          source,
		Destination:     dest,
		Package:         pkg,
		Imports:         imp,
		AuxFiles:        af,
		MockNames:       mn,
		SelfPackage:     spkg,
		CopyrightFile:   cf,
		WritePkgComment: wpcp,
		DebugParser:     dpp,
	}
}

func (r *Runner) SetSource(new string) {
	r.Source = new
}

func (r *Runner) SetDestination(new string) {
	r.Destination = new
}

func (r *Runner) GetDestination() string {
	return r.Destination
}

func (r *Runner) GetSource() string {
	return r.Source
}

func (r *Runner) String() string {
	return exec.Command("mockgen", r.options()...).String()
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
	if r.WritePkgComment != nil {
		opts = append(opts, "-write_package_comment="+strconv.FormatBool(*r.WritePkgComment))
	}
	if r.DebugParser != nil {
		opts = append(opts, "-debug_parser="+strconv.FormatBool(*r.DebugParser))
	}

	return opts
}
