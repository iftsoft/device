package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")
	appPar := config.GetAppParams()
	err, appCfg := config.GetAppConfig(appPar)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		core.StartFileLogger(appCfg.Logger)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start device application")
	log.Info(appPar.String())
	log.Info(appCfg.String())

	//err = linker.GetLinkerPorts(log)

	//srv := duplex.NewDuplexServer(appCfg.Duplex, log)
	//hnd := handler.NewHandlerManager(appCfg.Handlers)
	////hnd.RegisterReflexFactory(plugin.GetValidatorCheckerFactory())
	//hnd.SetupDuplexServer(srv)
	//srv.SetClientManager(hnd)
	//err = srv.StartListen()
	//if err == nil {
	//
	//	core.WaitForSignal(log)
	//
	//	srv.StopListen()
	//}
	log.Info("Stop device application")
	//hnd.Cleanup()
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
