package realmain

import (
	"log"
	"os"

	"github.com/sanposhiho/gomockhandler/realmain/util"

	"github.com/sanposhiho/gomockhandler/model"
)

func (r Runner) Generate() {
	if err := r.MockgenRunner.Run(); err != nil {
		log.Fatalf("failed to run mockgen: %v", err)
	}

	if r.Args.Destination != "" {
		chunk, err := r.ChunkRepo.Get()
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("failed to get chunk: %v", err)
			}
			chunk = model.NewChunk()
		}

		checksum, err := util.MockChackSum(r.Args.Destination)
		if err != nil {
			log.Fatalf("failed to calculate checksum of the mock: %v", err)
		}

		mock := model.NewMock(r.Args.Source, r.Args.Destination, checksum)
		chunk.PutMock(mock)
		if err := r.ChunkRepo.Put(chunk); err != nil {
			log.Fatalf("failed to put chunk: %v", err)
		}
	}
	return
}
