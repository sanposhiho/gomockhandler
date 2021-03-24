package command

import (
	"context"
	"fmt"
	"github.com/sanposhiho/gomockhandler/internal/mockgen"
	"log"
	"os"
	"path/filepath"

	"github.com/sanposhiho/gomockhandler/internal/model"
	"github.com/sanposhiho/gomockhandler/internal/util"
	"golang.org/x/sync/errgroup"
)

func (r Runner) Mockgen() {
	ch, err := r.ConfigRepo.Get(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	configPath, err := filepath.Abs(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get absolute project root: %v", err)
	}
	configDir, configFile := filepath.Split(configPath)
	if err := os.Chdir(configDir); err != nil {
		log.Fatalf("failed to change dir: %v", err)
	}

	sem := make(chan struct{}, r.Args.Concurrency)
	g, _ := errgroup.WithContext(context.Background())
	for _, m := range ch.Mocks {
		m := m
		sem <- struct{}{}

		g.Go(func() error {
			defer func() { <-sem }()
			var runner mockgen.Runner
			switch m.Mode {
			case model.ReflectMode:
				runner = m.ReflectModeRunner
			case model.SourceMode:
				runner = m.SourceModeRunner
			default:
				log.Printf("[WARN] unknown mock detected\n")
				return nil
			}
			err = runner.Run()
			if err != nil {
				return fmt.Errorf("run mockgen: %v", err)
			}

			checksum, err := util.MockCheckSum(runner.GetDestination())
			if err != nil {
				return fmt.Errorf("calculate checksum of the mock: %v", err)
			}

			m.CheckSum = checksum
			return nil
		})
	}
	err = g.Wait()
	close(sem)
	if err != nil {
		log.Fatalf("failed to run: %v", err.Error())
	}

	if err := r.ConfigRepo.Put(ch, configFile); err != nil {
		log.Fatalf("failed to put config: %v", err)
	}
	return
}
