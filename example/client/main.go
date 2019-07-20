package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy/system"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	appPar := config.GetAppParams()
	err, appCfg := config.GetAppConfig(appPar)
	if err != nil {
		fmt.Println(err)
	} else {
		core.StartFileLogger(&appCfg.Logger)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start application")
	log.Info("App params: %+v", appPar)
	if appCfg != nil {
		log.Info("Config logger: %+v", appCfg.Logger)
		log.Info("Config client: %+v", appCfg.Client)
		log.Info("Config device: %+v", appCfg.Device)
	}
	cln := duplex.NewDuplexClient(&appCfg.Client)

	logger := core.GetLogAgent(core.LogLevelTrace, "System")
	sysClnt := system.NewSystemClient()
	sysStub := system.NewSystemStub()
	sysClnt.Init(cln, sysStub, logger)
	sysStub.Init(sysClnt, logger)
	cln.AddScopeItem(sysClnt.GetScopeItem())

	cln.Start()

	WaitForSignal(log)

	cln.Stop()
	log.Info("Stop application")
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}

func WaitForSignal(out *core.LogAgent) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	out.Info("Got signal: %v, exiting.", s)
}
