package realmain

import (
	"log"
	"os"

	"github.com/sanposhiho/gomockhandler/model"
	"github.com/sanposhiho/gomockhandler/realmain/util"
)

func (r Runner) Check(tmpDir string) {
	if err := r.MockgenRunner.Run(); err != nil {
		log.Fatalf("failed to run mockgen: %v", err)
	}

	chunk, err := r.ChunkRepo.Get()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to get chunk: %v", err)
		}
		chunk = model.NewChunk()
	}

	checksum, err := util.MockChackSum(tmpDir)
	if err != nil {
		log.Fatalf("failed to calculate checksum of the mock: %v", err)
	}

	m, err := chunk.Find(r.Args.Destination)
	if err != nil {
		log.Fatalf("failed to get chunk: %v", err)
	}

	if m.CheckSum != checksum {
		// mock is not up to date
		log.Printf("[WARN] mock is not up to date. source: %s, destination: %s", r.Args.Source, r.Args.Destination)
	}
}
