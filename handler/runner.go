package handler

import (
	"errors"
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"os"
	"path/filepath"
	"sync"
)

type BinaryRunner struct {
	appArgs  []string 			// List of application params
	log      *core.LogAgent
	exitLoop bool
}

func newBinaryRunner(cfg *config.CommandConfig, log *core.LogAgent) *BinaryRunner {
	br := BinaryRunner{
		log:       log,
		appArgs:   make([]string, 0),
		exitLoop:  false,
	}
	br.processConfig(cfg)
	return &br
}

func (br *BinaryRunner) processConfig(cfg *config.CommandConfig) {
	br.appArgs = append(br.appArgs, cfg.BinaryFile)
	if cfg.DeviceName != "" {
		nameStr := fmt.Sprintf(`-name=%s`, cfg.DeviceName)
		br.appArgs = append(br.appArgs, nameStr)
	}
	if cfg.ConfigFile != "" {
		fullCfg, err := filepath.Abs(cfg.ConfigFile)
		if err == nil {
			cfgStr := fmt.Sprintf(`-cfg="%s"`, fullCfg)
			br.appArgs = append(br.appArgs, cfgStr)
		}
	}
	if cfg.Database != "" {
		fullBase, err := filepath.Abs(cfg.Database)
		if err == nil {
			baseStr := fmt.Sprintf(`-base="%s"`, fullBase)
			br.appArgs = append(br.appArgs, baseStr)
		}
	}
	if cfg.LoggerPath != "" {
		fullLog, err := filepath.Abs(cfg.LoggerPath)
		if err == nil {
			logsStr := fmt.Sprintf(`-logs="%s"`, fullLog)
			br.appArgs = append(br.appArgs, logsStr)
		}
	}
}


func (br *BinaryRunner) stopRunnerLoop() {
	br.exitLoop = true
}


func (br *BinaryRunner) launchRunnerLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	br.log.Info("Runner loop for file:%s is started", br.appArgs[0])
	defer br.log.Info("Runner loop for file:%s is stopped", br.appArgs[0])
	count := 0
	for br.exitLoop == false {
		br.log.Debug("Attempt %d to launch file: %s", count, br.appArgs[0])
		err := br.startBinary()
		if err != nil {
			count++
			if count > 3 {
				break
			}
		}
	}
}

func (br *BinaryRunner) startBinary() error {
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	br.log.Debug("BinaryRunner.startBinary for: %v", br.appArgs)
	process, err := os.StartProcess(br.appArgs[0], br.appArgs, &procAttr)
	if err == nil {
		var state *os.ProcessState
		state, err = process.Wait()
		if err == nil {
			br.log.Debug("BinaryRunner.process.Wait return: %d", state.ExitCode())
			if state.ExitCode() > 0 {
				err = errors.New("client process abnormal termination")
			}
		}
		if err != nil {
			br.log.Error("BinaryRunner.process.Wait error: %s", err.Error())
		}
	} else {
		br.log.Error("BinaryRunner.StartProcess error: %s", err.Error())
	}
	return err
}



