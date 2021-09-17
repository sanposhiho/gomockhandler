package command

import (
	"log"
	"path/filepath"
)

// DeleteMock delete a mock from config
func (r Runner) DeleteMock() {
	configPath := r.Args.ConfigPath

	chunk, err := r.ConfigRepo.Get(configPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	chunk.DeleteMock(filepath.Clean(r.Args.Destination))
	if err := r.ConfigRepo.Put(chunk, configPath); err != nil {
		log.Fatalf("failed to put config: %v", err)
	}

	log.Println("[INFO] The mock has been successfully deleted from configuration.")
	log.Println("Please delete the mock file itself manually.")

	return
}
