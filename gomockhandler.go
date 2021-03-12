package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/sanposhiho/gomockhandler/mockgen"

	"github.com/sanposhiho/gomockhandler/realmain"

	"github.com/sanposhiho/gomockhandler/mockgen/reflectmode"
	"github.com/sanposhiho/gomockhandler/mockgen/sourcemode"
	mockrepo "github.com/sanposhiho/gomockhandler/repository/config"
)

var (
	projectRoot = flag.String("project_root", "", "(for generate config) A config file will be generated in project root.")
	configPath  = flag.String("config", "./gomockhandler.json", "(for mockgen/chack) The path to config file.")
	concurrency = flag.Int("concurrency", 1, "Number of processes in parallel.")

	// flags for mockgen
	source          = flag.String("source", "", "(source mode) Input Go source file; enables source mode.")
	destination     = flag.String("destination", "", "Output file; defaults to stdout.")
	mockNames       = flag.String("mock_names", "", "Comma-separated interfaceName=mockName pairs of explicit config names to use. Mock names default to 'Mock'+ interfaceName suffix.")
	packageOut      = flag.String("package", "", "Package of the generated code; defaults to the package of the input with a 'mock_' prefix.")
	selfPackage     = flag.String("self_package", "", "The full package import path for the generated code. The purpose of this flag is to prevent import cycles in the generated code by trying to include its own package. This can happen if the config's package is set to one of its inputs (usually the main one) and the output is stdio so mockgen cannot detect the final output package. Setting this flag will then tell mockgen which import to exclude.")
	writePkgComment = flag.Bool("write_package_comment", true, "Writes package documentation comment (godoc) if true.")
	copyrightFile   = flag.String("copyright_file", "", "Copyright file used to add copyright header")
	imports         = flag.String("imports", "", "(source mode) Comma-separated name=path pairs of explicit imports to use.")
	auxFiles        = flag.String("aux_files", "", "(source mode) Comma-separated pkg=path pairs of auxiliary Go source files.")
	execOnly        = flag.String("exec_only", "", "(reflect mode) If set, execute this reflection program.")
	buildFlags      = flag.String("build_flags", "", "(reflect mode) Additional flags for go build.")
	progOnly        = flag.Bool("prog_only", false, "(reflect mode) Only generate the reflection program; write it to stdout and exit.")
	debugParser     = flag.Bool("debug_parser", false, "Print out parser results only.")
)

func main() {
	flag.Parse()

	repo := mockrepo.NewRepository()
	rm := realmain.Runner{
		ChunkRepo: &repo,
		Args: realmain.Args{
			ProjectRoot:     *projectRoot,
			ConfigPath:      *configPath,
			Concurrency:     *concurrency,
			Source:          *source,
			Destination:     *destination,
			MockNames:       *packageOut,
			PackageOut:      *packageOut,
			SelfPackage:     *selfPackage,
			WritePkgComment: *writePkgComment,
			CopyrightFile:   *copyrightFile,
			Imports:         *imports,
			AuxFiles:        *auxFiles,
			ExecOnly:        *execOnly,
			BuildFlags:      *buildFlags,
			ProgOnly:        *progOnly,
			DebugParser:     *debugParser,
		},
	}

	var realmain func()
	switch flag.Arg(0) {
	case "mockgen":
		realmain = rm.Mockgen
	case "check":
		realmain = rm.Check
	case "deletemock":
		realmain = rm.DeleteMock
	default:
		rm.MockgenRunner = prepareMockgenRunner()
		realmain = rm.GenerateConfig
	}

	realmain()
}

func prepareMockgenRunner() mockgen.Runner {
	if *source == "" {
		// reflect mode
		if flag.NArg() != 2 {
			log.Fatal("Expected exactly two arguments")
		}
		packageName := flag.Arg(0)
		interfaces := flag.Arg(1)
		if packageName == "." {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatalf("Get current directory failed: %v", err)
			}
			packageName, err = packageNameOfDir(dir)
			if err != nil {
				log.Fatalf("Parse package name failed: %v", err)
			}
		}
		return reflectmode.NewRunner(packageName, interfaces, *source, *destination, *packageOut, *imports, *auxFiles, *buildFlags, *mockNames, *selfPackage, *copyrightFile, *execOnly, *progOnly, *writePkgComment, *debugParser)
	}

	// source mode
	return sourcemode.NewRunner(*source, *destination, *packageOut, *imports, *auxFiles, *mockNames, *selfPackage, *copyrightFile, *writePkgComment, *debugParser)
}

// Plundered from golang/mock/mockgen/parse.go.
// packageNameOfDir get package import path via dir
func packageNameOfDir(srcDir string) (string, error) {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	var goFilePath string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			goFilePath = file.Name()
			break
		}
	}
	if goFilePath == "" {
		return "", fmt.Errorf("go source file not found %s", srcDir)
	}

	packageImport, err := parsePackageImport(srcDir)
	if err != nil {
		return "", err
	}
	return packageImport, nil
}

// Plundered from golang/mock/mockgen/parse.go.
//
// parseImportPackage get package import path via source file
// an alternative implementation is to use:
// cfg := &packages.Config{Mode: packages.NeedName, Tests: true, Dir: srcDir}
// pkgs, err := packages.Load(cfg, "file="+source)
// However, it will call "go list" and slow down the performance
func parsePackageImport(srcDir string) (string, error) {
	moduleMode := os.Getenv("GO111MODULE")
	// trying to find the module
	if moduleMode != "off" {
		currentDir := srcDir
		for {
			dat, err := ioutil.ReadFile(filepath.Join(currentDir, "go.mod"))
			if os.IsNotExist(err) {
				if currentDir == filepath.Dir(currentDir) {
					// at the root
					break
				}
				currentDir = filepath.Dir(currentDir)
				continue
			} else if err != nil {
				return "", err
			}
			modulePath := modfile.ModulePath(dat)
			return filepath.ToSlash(filepath.Join(modulePath, strings.TrimPrefix(srcDir, currentDir))), nil
		}
	}
	// fall back to GOPATH mode
	goPaths := os.Getenv("GOPATH")
	if goPaths == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}
	goPathList := strings.Split(goPaths, string(os.PathListSeparator))
	for _, goPath := range goPathList {
		sourceRoot := filepath.Join(goPath, "src") + string(os.PathSeparator)
		if strings.HasPrefix(srcDir, sourceRoot) {
			return filepath.ToSlash(strings.TrimPrefix(srcDir, sourceRoot)), nil
		}
	}
	return "", errOutsideGoPath
}

var errOutsideGoPath = errors.New("source directory is outside GOPATH")
