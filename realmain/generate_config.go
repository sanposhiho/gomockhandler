package realmain

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sanposhiho/gomockhandler/model"
	"github.com/sanposhiho/gomockhandler/realmain/util"
)

// GenerateConfig generate config
func (r Runner) GenerateConfig() {
	configPath, err := filepath.Abs(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get absolute project root: %w", err)
	}
	configDir, _ := filepath.Split(configPath)
	chunk, err := r.ChunkRepo.Get(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to get config: %v", err)
		}
		chunk = model.NewChunk()
	}

	if err := r.MockgenRunner.Run(); err != nil {
		log.Fatalf("failed to run mockgen: %v", err)
	}

	if r.Args.Destination != "" {
		// calculate mock's check sum
		checksum, err := util.MockCheckSum(r.Args.Destination)
		if err != nil {
			log.Fatalf("failed to calculate checksum of the mock: %v", err)
		}
		currentPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get config: %v", err)
		}
		destinationPathInPro := util.PathInProject(configDir, currentPath+"/"+r.Args.Destination)
		r.MockgenRunner.SetDestination(destinationPathInPro)

		if r.Args.Source != "" {
			sourcePathInPro := util.PathInProject(configDir, currentPath+"/"+r.Args.Source)
			r.MockgenRunner.SetSource(sourcePathInPro)
		}
		// store into config
		mock := model.NewMock(destinationPathInPro, checksum, r.MockgenRunner)
		chunk.PutMock(mock)
		if err := r.ChunkRepo.Put(chunk, configPath); err != nil {
			log.Fatalf("failed to put config: %v", err)
		}
	}
	return
}
