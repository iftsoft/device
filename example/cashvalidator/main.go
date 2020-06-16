package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/dbase/dbvalid"
	"github.com/iftsoft/device/driver"
	"github.com/iftsoft/device/driver/validator"
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

	err = CheckValidatorDatabase(appCfg.Storage, appCfg.Duplex.DevName)
	if err == nil{
		dev := driver.NewSystemDevice(appCfg)
		drv := validator.NewValidatorDriver()
		err = dev.InitDevice(drv)
		if err == nil {
			dev.StartDeviceLoop()

			core.WaitForSignal(log)

			dev.StopDeviceLoop()
		} else {
			log.Error("Can't start device: %s", err)
		}
	} else {
		log.Error("Can't check database file: %s", err)
	}
	log.Info("SysStop application")
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}

func CheckValidatorDatabase(conf *dbase.StorageConfig, devName string) error {
	err := core.CheckOrCreateFile(conf.FileName)
	if err == nil{
		store := dbase.GetNewDBaseStore(conf)
		err = store.Open()
		if err == nil{
			dbval := dbvalid.NewDBaseValidator(store, devName)
			err = dbval.CreateAllTables()
			_ = store.Close()
		}
	}
	return err
}
