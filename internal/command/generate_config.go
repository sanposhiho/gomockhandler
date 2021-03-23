package command

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sanposhiho/gomockhandler/internal/model"
	"github.com/sanposhiho/gomockhandler/internal/util"
)

// GenerateConfig generate config
func (r Runner) GenerateConfig() {
	configPath, err := filepath.Abs(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get absolute project root: %v", err)
	}

	if r.Args.Destination != "" {
		currentPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get config: %v", err)
		}
		configDir, configFile := filepath.Split(configPath)

		if err := os.Chdir(configDir); err != nil {
			log.Fatalf("failed to change dir: %v", err)
		}

		chunk, err := r.ChunkRepo.Get(configFile)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("failed to get config: %v", err)
			}
			chunk = model.NewChunk()
		}

		destinationPathInPro := util.PathInProject(configDir, currentPath+"/"+r.Args.Destination)
		r.MockgenRunner.SetDestination(destinationPathInPro)

		if r.Args.Source != "" {
			sourcePathInPro := util.PathInProject(configDir, currentPath+"/"+r.Args.Source)
			r.MockgenRunner.SetSource(sourcePathInPro)
		}

		if err := r.MockgenRunner.Run(); err != nil {
			log.Fatalf("failed to run mockgen: %v", err)
		}

		// calculate mock's check sum
		checksum, err := util.MockCheckSum(r.MockgenRunner.GetDestination())
		if err != nil {
			log.Fatalf("failed to calculate checksum of the mock: %v", err)
		}

		// store into config
		mock := model.NewMock(destinationPathInPro, checksum, r.MockgenRunner)
		chunk.PutMock(mock)
		if err := r.ChunkRepo.Put(chunk, configFile); err != nil {
			log.Fatalf("failed to put config: %v", err)
		}
	}
	return
}
