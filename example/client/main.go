package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver"
	"github.com/iftsoft/device/system"
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
		return
	}
	core.StartFileLogger(&appCfg.Logger)
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start application")
	log.Info("App params: %+v", appPar)
	if appCfg != nil {
		log.Info("Config logger: %+v", appCfg.Logger)
		log.Info("Config client: %+v", appCfg.Client)
		log.Info("Config device: %+v", appCfg.Device)
	}

	dev := system.NewSystemDevice(appCfg)
	drv := driver.NewDummyDriver()
	err = dev.InitDevice(drv)
	if err == nil {
		dev.StartDevice()

		WaitForSignal(log)

		dev.StopDevice()
	} else {
		log.Error("Can't start device: %s", err)
	}
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
