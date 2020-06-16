package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/handler"
	"github.com/iftsoft/device/linker"
	"os"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	appPar := config.GetAppParams()
	err, appCfg := config.GetSrvConfig(appPar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		core.StartFileLogger(appCfg.Logger)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start server application")
	log.Info(appPar.String())
	log.Info(appCfg.String())

	err = linker.GetLinkerPorts(log)

	srv := duplex.NewDuplexServer(appCfg.Duplex, log)
	hnd := handler.NewHandlerManager(appCfg.Handlers)
//	hnd.RegisterReflexFactory(handler.GetDeviceTesterFactory())
	hnd.SetupDuplexServer(srv)
	srv.SetClientManager(hnd)
	err = srv.StartListen()
	if err == nil {
		hnd.LaunchAllBinaries()

		core.WaitForSignal(log)

		hnd.StopAllBinaries()
		srv.StopListen()
	}
	log.Info("Stop server application")
	hnd.Cleanup()
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
