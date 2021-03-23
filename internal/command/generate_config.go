package command

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"

	"github.com/sanposhiho/gomockhandler/internal/model"
	"github.com/sanposhiho/gomockhandler/internal/util"
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

	if err := mockgen.TestRun(r.MockgenRunner); err != nil {
		log.Fatalf("failed to run mockgen: %v", err)
	}

	if r.Args.Destination != "" {
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
		mock := model.NewMock(destinationPathInPro, [16]byte{}, r.MockgenRunner)
		chunk.PutMock(mock)
		if err := r.ChunkRepo.Put(chunk, configPath); err != nil {
			log.Fatalf("failed to put config: %v", err)
		}
	}
	return
}
