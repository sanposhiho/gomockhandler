package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"

	"golang.org/x/sync/errgroup"

	"github.com/sanposhiho/gomockhandler/internal/model"
)

func (r Runner) Check() {
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

	isFail := false
	g, _ := errgroup.WithContext(context.Background())
	for _, m := range ch.Mocks {
		m := m

		var runner mockgen.Runner
		g.Go(func() error {
			switch m.Mode {
			case model.ReflectMode:
				runner = m.ReflectModeRunner
			case model.SourceMode:
				runner = m.SourceModeRunner
				sourceChecksum, err := mockgen.SourceChecksum(runner)
				if err != nil {
					return fmt.Errorf("failed to calculate checksum of the source: %v", err)
				}
				if sourceChecksum == m.SourceChecksum && !r.Args.ForceGenerate {
					// source file is not updated, so the mock is up-to-date.
					return nil
				}
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
	if isFail {
		log.Fatal("mocks is not up-to-date")
	}
	log.Print("mocks is up-to-date ✨")
	return
}
