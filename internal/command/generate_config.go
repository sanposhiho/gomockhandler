package command

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"

	"github.com/sanposhiho/gomockhandler/internal/model"
	"github.com/sanposhiho/gomockhandler/internal/util"
)

// GenerateConfig generates config
// require destination option.
func (r Runner) GenerateConfig() {
	configPath, err := filepath.Abs(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get absolute project root: %v", err)
	}

	// get the path where command is executed
	originalPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}
	configDir, configFile := filepath.Split(configPath)

	// move to the config directory
	if err := os.Chdir(configDir); err != nil {
		log.Fatalf("failed to change dir: %v", err)
	}

	// get config
	chunk, err := r.ConfigRepo.Get(configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to get config: %v", err)
		}
		chunk = model.NewChunk()
	}

	// change destination as seen from the config directory.
	destinationPathInPro := util.PathInProject(configDir, originalPath+"/"+r.Args.Destination)
	r.MockgenRunner.SetDestination(destinationPathInPro)

	var sourceChecksum string
	if r.Args.Source != "" {
		// change source as seen from the config directory.
		sourcePathInPro := util.PathInProject(configDir, originalPath+"/"+r.Args.Source)
		r.MockgenRunner.SetSource(sourcePathInPro)

		sourceChecksum, err = mockgen.SourceChecksum(r.MockgenRunner)
		if err != nil {
			log.Fatalf("failed to calculate checksum of the source: %v", err)
		}
	}

	if err := r.MockgenRunner.Run(); err != nil {
		log.Fatalf("failed to run mockgen: %v \nPlease run `%s` and check if mockgen works correctly with your options", err, r.MockgenRunner)
	}

	// calculate mock's check sum
	checksum, err := util.CalculateCheckSum(r.MockgenRunner.GetDestination())
	if err != nil {
		log.Fatalf("failed to calculate checksum of the mock: %v", err)
	}

	// store into config
	mock := model.NewMock(checksum, sourceChecksum, r.MockgenRunner)
	chunk.PutMock(r.MockgenRunner.GetDestination(), mock)
	if err := r.ConfigRepo.Put(chunk, configFile); err != nil {
		log.Fatalf("failed to put config: %v", err)
	}
	return
}
