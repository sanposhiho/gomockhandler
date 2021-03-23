package command

import (
	"context"
	"fmt"
	"log"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"

	"golang.org/x/sync/errgroup"

	"github.com/sanposhiho/gomockhandler/internal/model"
)

func (r Runner) Check() {
	ch, err := r.ChunkRepo.Get(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	isFail := false
	sem := make(chan struct{}, r.Args.Concurrency)
	g, _ := errgroup.WithContext(context.Background())
	for _, m := range ch.Mocks {
		m := m
		sem <- struct{}{}

		var runner mockgen.Runner
		g.Go(func() error {
			defer func() { <-sem }()

			switch m.Mode {
			case model.ReflectMode:
				runner = m.ReflectModeRunner
			case model.SourceMode:
				runner = m.SourceModeRunner
			default:
				log.Printf("[WARN] unknown mock detected\nPlease reconfigure the mock. destination: %s", runner.GetDestination())
				return nil
			}

			checksum, err := mockgen.Checksum(runner)
			if err != nil {
				return fmt.Errorf("get checksum: %v", err)
			}

			if m.CheckSum != checksum {
				// mock is not up to date
				s := runner.GetSource()
				d := runner.GetDestination()
				if s == "" {
					log.Printf("[ERROR] mock is not up to date. destination: %s", d)
				} else {
					log.Printf("[ERROR] mock is not up to date. source: %s, destination: %s", s, d)
				}
				isFail = true
			}
			return nil
		})
	}
	err = g.Wait()
	close(sem)
	if err != nil {
		log.Fatalf("failed to run: %v", err.Error())
	}
	if isFail {
		log.Fatal("mocks is not up-to-date")
	}
	return
}
