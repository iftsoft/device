package main

import (
	"fmt"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/dbase/dbvalid"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	appPar := config.GetAppParams()
	devCfg := config.GetDefaultDeviceConfig()
	err, appCfg := config.GetAppConfig(appPar, devCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	core.StartFileLogger(appCfg.Logger)
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("SysStart application")
	log.Info(appPar.String())
	log.Info(appCfg.String())

	err = core.CheckOrCreateFile(appCfg.Storage.FileName)
	if err != nil{
		log.Error("Can't check dbase file: %s", err)
	} else {
		store := dbase.GetNewDBaseStore(appCfg.Storage)
		err = store.Open()
		if err != nil{
			log.Error("Can't open dbase file: %s", err)
		} else {
			dbval := dbvalid.NewDBaseValidator(store, appCfg.Duplex.DevName)
			err = dbval.CreateAllTables()
			if err != nil{
				log.Error("Can't create tables: %s", err)
			}
			batch := &dbvalid.ObjBatch{}
			batch.Device = "cashvalidator"
			err = dbval.GetLastBatch(batch)
			if err != nil{
				log.Error("Can't read batch: %s", err)
			}
			log.Info("Batch ID: %d, device: %s, time: %v", batch.Id, batch.Device, batch.Opened)
			defNoteList := common.ValidNoteList {
				{appCfg.Duplex.DevName, 980, 10, 1.0, 10.0, },
				{appCfg.Duplex.DevName, 980, 10, 2.0, 20.0, },
				{appCfg.Duplex.DevName, 980, 10, 5.0, 50.0, },
			}
			err = dbval.InitNoteList(defNoteList)
			if err != nil{
				log.Error("Can't init notes: %s", err)
			}
			list, err := dbval.ReadNoteList()
			if err != nil{
				log.Error("Can't init notes: %s", err)
			}
			log.Info(list.String())
		}
		err = store.Close()
	}

//	err = linker.GetLinkerPorts(log)

	//dev := driver.NewSystemDevice(appCfg)
	//drv := loopback.NewDummyDriver()
	//err = dev.InitDevice(drv)
	//if err == nil {
	//	dev.StartDevice()
	//
	//	core.WaitForSignal(log)
	//
	//	dev.StopDevice()
	//} else {
	//	log.Error("Can't start device: %s", err)
	//}
	log.Info("SysStop application")
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
