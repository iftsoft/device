package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver"
	"github.com/iftsoft/device/system"
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
	core.StartFileLogger(appCfg.Logger)
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start application")
	log.Info(appPar.String())
	log.Info(appCfg.String())

	dev := system.NewSystemDevice(appCfg)
	drv := driver.NewDummyDriver()
	err = dev.InitDevice(drv)
	if err == nil {
		dev.StartDevice()

		core.WaitForSignal(log)

		dev.StopDevice()
	} else {
		log.Error("Can't start device: %s", err)
	}
	log.Info("Stop application")
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
