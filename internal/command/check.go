package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/semaphore"

	"golang.org/x/sync/errgroup"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"
	"github.com/sanposhiho/gomockhandler/internal/model"
	"github.com/sanposhiho/gomockhandler/internal/zombie"
)

func (r Runner) Check() {
	ctx := context.Background()
	ch, err := r.ConfigRepo.Get(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	configPath, err := filepath.Abs(r.Args.ConfigPath)
	if err != nil {
		log.Fatalf("failed to get absolute project root: %v", err)
	}
	configDir, _ := filepath.Split(configPath)
	if err := os.Chdir(configDir); err != nil {
		log.Fatalf("failed to change dir: %v", err)
	}

	zombieChecker := zombie.NewChecker()
	if err = zombieChecker.FindMocks(ctx); err != nil {
		log.Fatalf("failed to find mock list generated by gomock: %v\n", err)
	}

	isFail := false
	g, _ := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	for _, m := range ch.Mocks {
		m := m
		sem.Acquire(ctx, 1)

		var runner mockgen.Runner
		g.Go(func() error {
			defer sem.Release(1)
			switch m.Mode {
			case model.ReflectMode:
				runner = m.ReflectModeRunner
			case model.SourceMode:
				runner = m.SourceModeRunner
			default:
				log.Printf("[WARN] unknown mock detected\nPlease reconfigure the mock. destination: %s", runner.GetDestination())
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

			key := runner.GetDestination()
			if ok := zombieChecker.Search(key); !ok {
				fmt.Fprintf(os.Stderr, "[ERROR] mock has not been generated. destination: %s\n", key)
				isFail = true
				return nil
			}

			if m.Mode == model.SourceMode {
				sourceChecksum, err := mockgen.SourceChecksum(runner)
				if err != nil {
					return fmt.Errorf("failed to calculate checksum of the source: %v", err)
				}
				if sourceChecksum == m.SourceChecksum && !r.Args.ForceGenerate {
					// source file is not updated, so the mock is up-to-date.
					return nil
				}
			}

			checksum, err := mockgen.Checksum(runner)
			if err != nil {
				return fmt.Errorf("get checksum: %v", err)
			}

			if m.MockCheckSum != checksum {
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
	if err != nil {
		log.Fatalf("failed to run: %v", err.Error())
	}

	if zombieList := zombieChecker.Check(); len(zombieList) != 0 {
		fmt.Fprintf(os.Stderr, "[WARN] Some mocks are un-managed by gomockhandler: %v\n", zombieList)
		isFail = true
	}
	if isFail {
		log.Fatal("mocks is not up-to-date")
	}
	log.Print("mocks is up-to-date ✨")

	return
}
