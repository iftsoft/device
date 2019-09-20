package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/handler"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	appPar := config.GetAppParams()
	err, appCfg := config.GetSrvConfig(appPar)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		core.StartFileLogger(appCfg.Logger)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start server application")
	log.Info(appPar.String())
	log.Info(appCfg.String())

	srv := duplex.NewDuplexServer(appCfg.Duplex, log)
	obj := handler.NewObjectProxy()
	obj.Init(srv)
	srv.SetClientManager(obj.GetClientManager())
	err = srv.StartListen()
	if err == nil {

		core.WaitForSignal(log)

		srv.StopListen()
	}
	log.Info("Stop server application")
	obj.Cleanup()
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
