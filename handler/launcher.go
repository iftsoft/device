package handler

import (
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"sync"
)


type BinaryLauncher struct {
	runnerList []*BinaryRunner
	log        *core.LogAgent
	wg         sync.WaitGroup
}

func (bl *BinaryLauncher) initLauncher(list config.HandlerList) {
	bl.runnerList = make([]*BinaryRunner,0)
	bl.log = core.GetLogAgent(core.LogLevelTrace, "Runner")
	for _, cfg := range list {
		if cfg.Command.Enabled == true &&
			cfg.Command.BinaryFile != "" {
			runner := newBinaryRunner(&cfg.Command, bl.log)
			bl.runnerList = append(bl.runnerList, runner)
			bl.log.Trace("BinaryLauncher.append runner for %s", cfg.Command.BinaryFile)
		}
	}
}

func (bl *BinaryLauncher) launchAllBinaries() {
	bl.log.Trace("BinaryLauncher.launchAllBinaries")
	for _, run := range bl.runnerList {
		go run.launchRunnerLoop(&bl.wg)
	}
}

func (bl *BinaryLauncher) setQuitFlag() {
	bl.log.Trace("BinaryLauncher.setQuitFlag")
	for _, run := range bl.runnerList {
		run.stopRunnerLoop()
	}
}

func (bl *BinaryLauncher) waitAll() {
	bl.log.Trace("BinaryLauncher.waitAll")
	bl.wg.Wait()
}

