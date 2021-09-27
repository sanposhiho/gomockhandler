package runnermanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner32"

	"golang.org/x/sync/semaphore"

	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner10"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner11"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner12"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner13"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner14"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner15"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner16"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner17"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner18"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner19"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner2"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner20"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner21"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner22"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner23"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner24"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner25"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner26"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner27"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner28"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner29"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner3"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner30"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner31"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner4"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner5"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner6"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner7"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner8"
	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner9"

	"github.com/sanposhiho/gomockhandler/internal/mockgen"

	"github.com/sanposhiho/gomockhandler/internal/mockgen/native/runners/runner1"
)

type nativeRunner interface {
	mockgen.Runner

	Set(pn, ifs, source, dest, pkg, imp, af, bf, mn, spkg, cf, eo string, po, wpc, dp *bool)
}

type Manager struct {
	runners map[int]nativeRunner
	isUsed  map[int]bool
	sem     *semaphore.Weighted
	mu      sync.Mutex
}

const (
	runnerNum = 32
)

func New(concurrency int64) *Manager {
	if concurrency >= runnerNum {
		concurrency = runnerNum
	}

	mu := &sync.Mutex{}

	m := &Manager{
		runners: map[int]nativeRunner{
			0:  runner32.NewRunner(mu),
			1:  runner1.NewRunner(mu),
			2:  runner2.NewRunner(mu),
			3:  runner3.NewRunner(mu),
			4:  runner4.NewRunner(mu),
			5:  runner5.NewRunner(mu),
			6:  runner6.NewRunner(mu),
			7:  runner7.NewRunner(mu),
			8:  runner8.NewRunner(mu),
			9:  runner9.NewRunner(mu),
			10: runner10.NewRunner(mu),
			11: runner11.NewRunner(mu),
			12: runner12.NewRunner(mu),
			13: runner13.NewRunner(mu),
			14: runner14.NewRunner(mu),
			15: runner15.NewRunner(mu),
			16: runner16.NewRunner(mu),
			17: runner17.NewRunner(mu),
			18: runner18.NewRunner(mu),
			19: runner19.NewRunner(mu),
			20: runner20.NewRunner(mu),
			21: runner21.NewRunner(mu),
			22: runner22.NewRunner(mu),
			23: runner23.NewRunner(mu),
			24: runner24.NewRunner(mu),
			25: runner25.NewRunner(mu),
			26: runner26.NewRunner(mu),
			27: runner27.NewRunner(mu),
			28: runner28.NewRunner(mu),
			29: runner29.NewRunner(mu),
			30: runner30.NewRunner(mu),
			31: runner31.NewRunner(mu),
		},
		isUsed: map[int]bool{},
	}

	for i := 0; int64(i) < concurrency; i++ {
		m.isUsed[i] = false
	}

	m.sem = semaphore.NewWeighted(concurrency)

	return m
}

func (m *Manager) Run(pn, ifs, source, dest, pkg, imp, af, bf, mn, spkg, cf, eo string, po, wpc, dp *bool) error {
	err := m.sem.Acquire(context.Background(), 1)
	if err != nil {
		return fmt.Errorf("acquire semaphore: %w", err)
	}
	defer m.sem.Release(1)

	for i, v := range m.isUsed {
		if !v {
			m.isUsed[i] = true
			defer func() {
				m.isUsed[i] = false
			}()

			m.runners[i].Set(pn, ifs, source, dest, pkg, imp, af, bf, mn, spkg, cf, eo, po, wpc, dp)
			if err := m.runners[i].Run(); err != nil {
				return err
			}

			return nil
		}
	}
	return nil
}
