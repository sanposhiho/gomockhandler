package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	if err != nil {
		log.Fatalf("failed to run: %v", err.Error())
	}
	if isFail {
		log.Fatal("mocks is not up-to-date")
	}
	return
}
