package reflectmode

import (
	"os/exec"
	"strconv"
)

type Runner struct {
	PackageName     string `json:"package_name"`
	Interfaces      string `json:"interfaces"`
	Source          string `json:"source"`
	Destination     string `json:"destination"`
	Package         string `json:"package"`
	Imports         string `json:"imports"`
	AuxFiles        string `json:"aux_files"`
	BuildFlags      string `json:"build_flags"`
	MockNames       string `json:"mock_names"`
	SelfPackage     string `json:"self_package"`
	CopyrightFile   string `json:"copyright_file"`
	ExecOnly        string `json:"exec_only"`
	ProgOnly        bool   `json:"prog_only"`
	WritePkgComment bool   `json:"write_pkg_comment"`
	DebugParser     bool   `json:"debug_parser"`
}

func NewRunner(pn, ifs, source, dest, pkg, imp, af, bf, mn, spkg, cf, eo string, po, wpc, dp bool) *Runner {
	return &Runner{
		PackageName:     pn,
		Interfaces:      ifs,
		Source:          source,
		Destination:     dest,
		Package:         pkg,
		Imports:         imp,
		AuxFiles:        af,
		BuildFlags:      bf,
		MockNames:       mn,
		SelfPackage:     spkg,
		CopyrightFile:   cf,
		ExecOnly:        eo,
		ProgOnly:        po,
		WritePkgComment: wpc,
		DebugParser:     dp,
	}
}

func (r *Runner) SetSource(new string) {
	r.Source = new
}

func (r *Runner) SetDestination(new string) {
	r.Destination = new
}

func (r *Runner) Run() error {
	params := append(r.options(), []string{r.PackageName, r.Interfaces}...)
	return exec.Command("mockgen", params...).Run()
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
	if r.BuildFlags != "" {
		opts = append(opts, "-build_flags="+r.BuildFlags)
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
	if r.ExecOnly != "" {
		opts = append(opts, "-exec_only="+r.ExecOnly)
	}
	opts = append(opts, "-prog_only="+strconv.FormatBool(r.ProgOnly))
	opts = append(opts, "-write_package_comment="+strconv.FormatBool(r.WritePkgComment))
	opts = append(opts, "-debug_parser="+strconv.FormatBool(r.DebugParser))

	return opts
}
