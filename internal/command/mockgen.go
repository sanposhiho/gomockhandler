package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"

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

	g, _ := errgroup.WithContext(context.Background())
	for _, m := range ch.Mocks {
		m := m
		g.Go(func() error {
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
				return fmt.Errorf("failed to run mockgen: %v \nPlease run `%s` and check if mockgen works correctly with your options", err, r.MockgenRunner)
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
	if err != nil {
		log.Fatalf("failed to run: %v", err.Error())
	}

	if err := r.ConfigRepo.Put(ch, configFile); err != nil {
		log.Fatalf("failed to put config: %v", err)
	}
	return
}
