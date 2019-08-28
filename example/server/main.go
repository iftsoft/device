package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/handler"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	appPar := config.GetSrvParams()
	err, appCfg := config.GetSrvConfig(appPar)
	if err != nil {
		fmt.Println(err)
	} else {
		core.StartFileLogger(&appCfg.Logger)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start server application")
	log.Info("Server app params: %+v", appPar)
	if appCfg != nil {
		log.Info("Config logger: %+v", appCfg.Logger)
		log.Info("Config server: %+v", appCfg.Server)
	}
	srv := duplex.NewDuplexServer(&appCfg.Server, log)
	obj := handler.NewObjectProxy()
	obj.Init(srv)
	srv.SetClientManager(obj)
	err = srv.StartListen()
	if err == nil {

		WaitForSignal(log)

		srv.StopListen()
	}
	log.Info("Stop server application")
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
