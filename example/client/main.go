package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver"
	"github.com/iftsoft/device/driver/loopback"
	"github.com/iftsoft/device/linker"
	"os"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	appPar := config.GetAppParams()
	devCfg := config.GetDefaultDeviceConfig()
	err, appCfg := config.GetAppConfig(appPar, devCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	core.StartFileLogger(appCfg.Logger)
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("SysStart application")
	log.Info(appPar.String())
	log.Info(appCfg.String())

	err = linker.GetLinkerPorts(log)

	dev := driver.NewSystemDevice(appCfg)
	drv := loopback.NewDummyDriver()
	err = dev.InitDevice(drv)
	if err == nil {
		dev.StartDeviceLoop()

		core.WaitForSignal(log)

		dev.StopDeviceLoop()
	} else {
		log.Error("Can't start device: %s", err)
	}
	log.Info("SysStop application")
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
