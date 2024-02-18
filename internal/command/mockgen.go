package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"
	"github.com/sanposhiho/gomockhandler/internal/model"
	"github.com/sanposhiho/gomockhandler/internal/util"
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
	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	for _, m := range ch.Mocks {
		m := m
		sem.Acquire(context.Background(), 1)

		g.Go(func() error {
			defer sem.Release(1)
			var runner mockgen.Runner
			var sourceChecksum string
			switch m.Mode {
			case model.ReflectMode:
				runner = m.ReflectModeRunner
			case model.SourceMode:
				runner = m.SourceModeRunner
				sourceChecksum, err = mockgen.SourceChecksum(runner)
				if err != nil {
					return fmt.Errorf("failed to calculate checksum of the source: %v", err)
				}
				if sourceChecksum == m.SourceChecksum && !r.Args.ForceGenerate {
					_, err := os.Stat(runner.GetDestination())
					if err == nil {
						// source file is not updated and mock file is exist, so we don't have to generate mock.
						return nil
					}
				}
			default:
				log.Printf("[WARN] unknown mock detected\n")
				return nil
			}
			if r.Args.PathFilter != "" {
				dest, err := filepath.Abs(runner.GetDestination())
				if err != nil {
					return fmt.Errorf("failed to get absolute path from mock's destination, please make sure the destination is correct, destination: %s, :%w", runner.GetDestination(), err)
				}
				pf, err := filepath.Abs(r.Args.PathFilter)
				if err != nil {
					return fmt.Errorf("failed to get absolute path from your filter option, please make sure the filter option is correct, filter option: %s, :%w", r.Args.PathFilter, err)
				}
				if !strings.HasPrefix(dest, pf+"/") {
					// skip
					return nil
				}
			}

			err = runner.Run()
			if err != nil {
				return fmt.Errorf("failed to run mockgen: %v \nPlease run `%s` and check if mockgen works correctly with your options", err, runner)
			}
			file, err := ioutil.ReadFile(r.MockgenRunner.GetDestination())
			if err != nil {
				log.Fatalf("failed read file. filename: %s, err: %w", r.MockgenRunner.GetDestination(), err)
			}
			checksum, err := util.CalculateCheckSum(file)
			if err != nil {
				return fmt.Errorf("calculate checksum of the mock: %v", err)
			}

			m.MockCheckSum = checksum
			m.SourceChecksum = sourceChecksum
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
