package realmain

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/sanposhiho/gomockhandler/model"
	"github.com/sanposhiho/gomockhandler/realmain/util"
)

func (r Runner) Check() {
	ch, err := r.ChunkRepo.Get(r.Args.ConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to get config: %v", err)
		}
	}

	// TODO: run in parallels
	for _, m := range ch.Mocks {
		d, err := ioutil.TempDir(".", "gomockhandler")
		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.RemoveAll(d)
		tmpFile := d + "/tmpmock.go"
		var source string
		var destination string
		switch m.Mode {
		case model.Unknown:
			log.Printf("unknown mock detected\n")
			continue
		case model.ReflectMode:
			source = m.ReflectModeRunner.Source
			destination = m.ReflectModeRunner.Destination
			m.ReflectModeRunner.SetDestination(tmpFile)
			err = m.ReflectModeRunner.Run()
		case model.SourceMode:
			source = m.SourceModeRunner.Source
			destination = m.SourceModeRunner.Destination
			m.SourceModeRunner.SetDestination(tmpFile)
			err = m.SourceModeRunner.Run()
		}
		if err != nil {
			log.Fatalf("failed to run mockgen: %v", err)
		}

		checksum, err := util.MockChackSum(tmpFile)
		if err != nil {
			log.Fatalf("failed to calculate checksum of the mock: %v", err)
		}

		if m.CheckSum != checksum {
			// mock is not up to date
			log.Printf("[WARN] mock is not up to date. source: %s, destination: %s", source, destination)
		}
	}
	return
}
