package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sanposhiho/gomockhandler/mockgen/reflectmode"

	"github.com/sanposhiho/gomockhandler/mockgen/sourcemode"

	mockrepo "github.com/sanposhiho/gomockhandler/repository/chunk"

	"github.com/sanposhiho/gomockhandler/model"
)

var (
	check = flag.Bool("check", false, "If true, check if mock is up to date")

	// flags for mockgen
	source          = flag.String("source", "", "(source mode) Input Go source file; enables source mode.")
	destination     = flag.String("destination", "", "Output file; defaults to stdout.")
	mockNames       = flag.String("mock_names", "", "Comma-separated interfaceName=mockName pairs of explicit chunk names to use. Mock names default to 'Mock'+ interfaceName suffix.")
	packageOut      = flag.String("package", "", "Package of the generated code; defaults to the package of the input with a 'mock_' prefix.")
	selfPackage     = flag.String("self_package", "", "The full package import path for the generated code. The purpose of this flag is to prevent import cycles in the generated code by trying to include its own package. This can happen if the chunk's package is set to one of its inputs (usually the main one) and the output is stdio so mockgen cannot detect the final output package. Setting this flag will then tell mockgen which import to exclude.")
	writePkgComment = flag.Bool("write_package_comment", true, "Writes package documentation comment (godoc) if true.")
	copyrightFile   = flag.String("copyright_file", "", "Copyright file used to add copyright header")
	imports         = flag.String("imports", "", "(source mode) Comma-separated name=path pairs of explicit imports to use.")
	auxFiles        = flag.String("aux_files", "", "(source mode) Comma-separated pkg=path pairs of auxiliary Go source files.")
	execOnly        = flag.String("exec_only", "", "(reflect mode) If set, execute this reflection program.")
	buildFlags      = flag.String("build_flags", "", "(reflect mode) Additional flags for go build.")
	progOnly        = flag.Bool("prog_only", false, "(reflect mode) Only generate the reflection program; write it to stdout and exit.")
	debugParser     = flag.Bool("debug_parser", false, "Print out parser results only.")
)

type chunkRepo interface {
	Put(m *model.Chunk) error
	Get() (*model.Chunk, error)
}

type mockgenRunner interface {
	Run() error
}

type Realmain struct {
	chunkRepo     chunkRepo
	mockgenRunner mockgenRunner
}

func main() {
	flag.Parse()

	repo := mockrepo.NewRepository()
	rm := Realmain{
		chunkRepo: &repo,
	}

	var realmain func()
	if *check {
		if *destination == "" {
			log.Fatal("Need destination to check if the mock is up-to-date.")
		}

		d, err := ioutil.TempDir(".", "")
		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.RemoveAll(d)
		tmpFile := d + "/tmpmock.go"

		if *source == "" {
			if flag.NArg() != 2 {
				log.Fatal("Expected exactly two arguments")
			}
			packageName := flag.Arg(0)
			interfaces := flag.Arg(1)
			r := reflectmode.NewRunner(packageName, interfaces, *source, tmpFile, *packageOut, *imports, *auxFiles, *buildFlags, *mockNames, *selfPackage, *copyrightFile, *execOnly, *progOnly, *writePkgComment, *debugParser)
			rm.mockgenRunner = r
		} else {
			r := sourcemode.NewRunner(*source, tmpFile, *packageOut, *imports, *auxFiles, *mockNames, *selfPackage, *copyrightFile, *writePkgComment, *debugParser)
			rm.mockgenRunner = r
		}

		realmain = func() {
			rm.check(tmpFile)
		}
	} else {
		if *source == "" {
			if flag.NArg() != 2 {
				log.Fatal("Expected exactly two arguments")
			}
			packageName := flag.Arg(0)
			interfaces := flag.Arg(1)
			r := reflectmode.NewRunner(packageName, interfaces, *source, *destination, *packageOut, *imports, *auxFiles, *buildFlags, *mockNames, *selfPackage, *copyrightFile, *execOnly, *progOnly, *writePkgComment, *debugParser)
			rm.mockgenRunner = r
		} else {
			r := sourcemode.NewRunner(*source, *destination, *packageOut, *imports, *auxFiles, *mockNames, *selfPackage, *copyrightFile, *writePkgComment, *debugParser)
			rm.mockgenRunner = r
		}

		realmain = rm.generate
	}

	realmain()
}

func (r Realmain) check(tmpDir string) {
	if err := r.mockgenRunner.Run(); err != nil {
		log.Fatalf("failed to run mockgen: %v", err)
	}

	chunk, err := r.chunkRepo.Get()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to get chunk: %v", err)
		}
		chunk = model.NewChunk()
	}

	checksum, err := mockChackSum(tmpDir)
	if err != nil {
		log.Fatalf("failed to calculate checksum of the mock: %v", err)
	}

	m, err := chunk.Find(*destination)
	if err != nil {
		log.Fatalf("failed to get chunk: %v", err)
	}

	if m.CheckSum != checksum {
		// mock is not up to date
		log.Printf("[WARN] mock is not up to date. source: %s, destination: %s", *source, *destination)
	}
}

func (r Realmain) generate() {
	if err := r.mockgenRunner.Run(); err != nil {
		log.Fatalf("failed to run mockgen: %v", err)
	}

	if *destination != "" {
		chunk, err := r.chunkRepo.Get()
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("failed to get chunk: %v", err)
			}
			chunk = model.NewChunk()
		}

		checksum, err := mockChackSum(*destination)
		if err != nil {
			log.Fatalf("failed to calculate checksum of the mock: %v", err)
		}

		mock := model.NewMock(*source, *destination, checksum)
		chunk.PutMock(mock)
		if err := r.chunkRepo.Put(chunk); err != nil {
			log.Fatalf("failed to put chunk: %v", err)
		}
	}
	return
}

func mockChackSum(filePath string) ([16]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed read file. filename: %s, err: %w", filePath)
	}

	hash := md5.Sum(file)
	return hash, nil
}
