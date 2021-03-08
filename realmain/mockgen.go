package realmain

import (
	"log"
	"strings"

	"github.com/sanposhiho/gomockhandler/model"
	"github.com/sanposhiho/gomockhandler/realmain/util"
)

func (r Runner) Mockgen() {
	ch, err := r.ChunkRepo.Get(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	// TODO: run in parallels
	for _, m := range ch.Mocks {
		var destination string
		switch m.Mode {
		case model.Unknown:
			log.Printf("unknown mock detected\n")
			continue
		case model.ReflectMode:
			err = m.ReflectModeRunner.Run()
			destination = m.ReflectModeRunner.Destination
		case model.SourceMode:
			err = m.SourceModeRunner.Run()
			destination = m.SourceModeRunner.Destination
		}
		if err != nil {
			log.Fatalf("failed to run mockgen: %v", err)
		}

		checksum, err := util.MockChackSum(destination)
		if err != nil {
			log.Fatalf("failed to calculate checksum of the mock: %v", err)
		}

		m.CheckSum = checksum
	}

	if err := r.ChunkRepo.Put(ch, r.Args.ConfigPath); err != nil {
		log.Fatalf("failed to put config: %v", err)
	}
	return
}

func pathInProject(projectRoot, path string) string {
	return strings.Replace(path, projectRoot, ".", 1)
}
