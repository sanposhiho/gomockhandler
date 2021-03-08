package realmain

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

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

	sem := make(chan struct{}, r.Args.Concurrency)
	g, _ := errgroup.WithContext(context.Background())
	for _, m := range ch.Mocks {
		sem <- struct{}{}

		g.Go(func() error {
			defer func() { <-sem }()

			tmpFile := tmpFilePath(m.Destination)
			defer os.Remove(tmpFile)
			var source string
			var destination string
			switch m.Mode {
			case model.Unknown:
				log.Printf("unknown mock detected\n")
				return nil
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
				return fmt.Errorf("run mockgen: %v", err)
			}

			checksum, err := util.MockChackSum(tmpFile)
			if err != nil {
				return fmt.Errorf("calculate checksum of the mock: %v", err)
			}

			if m.CheckSum != checksum {
				// mock is not up to date
				log.Printf("[WARN] mock is not up to date. source: %s, destination: %s", source, destination)
			}
			return nil
		})
	}
	err = g.Wait()
	close(sem)
	if err != nil {
		log.Fatalf("failed to run: %v", err.Error())
	}
	return
}

func tmpFilePath(original string) string {
	d, f := filepath.Split(original)
	return d + "tmp_" + f
}
