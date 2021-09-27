package runnertemplate

import (
	"flag"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type Runner struct {
	mu *sync.Mutex

	PackageName     string `json:"package_name,omitempty"`
	Interfaces      string `json:"interfaces,omitempty"`
	Source          string `json:"source,omitempty"`
	Destination     string `json:"destination,omitempty"`
	Package         string `json:"package,omitempty"`
	Imports         string `json:"imports,omitempty"`
	AuxFiles        string `json:"aux_files,omitempty"`
	BuildFlags      string `json:"build_flags,omitempty"`
	MockNames       string `json:"mock_names,omitempty"`
	SelfPackage     string `json:"self_package,omitempty"`
	CopyrightFile   string `json:"copyright_file,omitempty"`
	ExecOnly        string `json:"exec_only,omitempty"`
	ProgOnly        *bool  `json:"prog_only,omitempty"`
	WritePkgComment *bool  `json:"write_pkg_comment,omitempty"`
	DebugParser     *bool  `json:"debug_parser,omitempty"`
}

func NewRunner(m *sync.Mutex) *Runner {
	return &Runner{mu: m}
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
	params := append(r.options(), []string{r.PackageName, r.Interfaces}...)
	return exec.Command("mockgen", params...).String()
}

func (r *Runner) Run() error {
	t := true
	f := false

	if r.WritePkgComment == nil {
		r.WritePkgComment = &t
	}
	if r.ProgOnly == nil {
		r.ProgOnly = &f
	}
	if r.DebugParser == nil {
		r.DebugParser = &f
	}

	source = &r.Source
	destination = &r.Destination
	packageOut = &r.Package
	imports = &r.Imports
	auxFiles = &r.AuxFiles
	buildFlags = &r.BuildFlags
	mockNames = &r.MockNames
	selfPackage = &r.SelfPackage
	copyrightFile = &r.CopyrightFile
	execOnly = &r.ExecOnly
	progOnly = r.ProgOnly
	writePkgComment = r.WritePkgComment
	debugParser = r.DebugParser

	showVersion = &f
	r.mu.Lock()
	// reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// set os.Args
	os.Args = []string{"mockgen", r.PackageName, r.Interfaces}

	go func() {
		for {
			// unlock after flag parse
			if flag.Parsed() {
				r.mu.Unlock()
				break
			}
		}
	}()

	main()

	return nil
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
	if r.WritePkgComment != nil {
		opts = append(opts, "-write_package_comment="+strconv.FormatBool(*r.WritePkgComment))
	}
	if r.DebugParser != nil {
		opts = append(opts, "-debug_parser="+strconv.FormatBool(*r.DebugParser))
	}
	if r.ProgOnly != nil {
		opts = append(opts, "-prog_only="+strconv.FormatBool(*r.ProgOnly))
	}
	return opts
}

func (r *Runner) Set(pn, ifs, source, dest, pkg, imp, af, bf, mn, spkg, cf, eo string, po, wpc, dp *bool) {
	r.PackageName = pn
	r.Interfaces = ifs
	r.Source = source
	r.Destination = dest
	r.Package = pkg
	r.Imports = imp
	r.AuxFiles = af
	r.BuildFlags = bf
	r.MockNames = mn
	r.SelfPackage = spkg
	r.CopyrightFile = cf
	r.ExecOnly = eo
	r.ProgOnly = po
	r.WritePkgComment = wpc
	r.DebugParser = dp
}
