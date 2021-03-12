package realmain

import (
	"log"
	"path/filepath"
)

// DeleteMock delete a mock from config
func (r Runner) DeleteMock() {
	configPath := r.Args.ConfigPath

	chunk, err := r.ChunkRepo.Get(configPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	chunk.DeleteMock(filepath.Clean(r.Args.Destination))
	if err := r.ChunkRepo.Put(chunk, configPath); err != nil {
		log.Fatalf("failed to put config: %v", err)
	}
	return
}
