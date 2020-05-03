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
			err = dbval.GetLastBatch(batch)
			if err != nil{
				log.Error("Can't read batch: %s", err)
			}
			log.Info("Get last batch - %s", batch.String())

			defNoteList := common.ValidNoteList {
				{appCfg.Duplex.DevName, 980, 0, 1.0, 0.0, },
				{appCfg.Duplex.DevName, 980, 0, 2.0, 0.0, },
				{appCfg.Duplex.DevName, 980, 0, 5.0, 0.0, },
			}
			data := &common.ValidatorBatch{}
			err = dbval.CloseBatch(data)
			if err != nil{
				log.Error("Can't close batch: %s", err)
			}
			log.Info("Close batch - %s", data.String())

			err = dbval.InitNoteList(defNoteList)
			if err != nil{
				log.Error("Can't init notes: %s", err)
			}

			err = dbval.GetLastBatch(batch)
			if err != nil{
				log.Error("Can't read batch: %s", err)
			}
			log.Info("Get last batch - %s", batch.String())

			accept1 := common.ValidatorAccept{ 980, 1.0,1, 1.0}
			accept2 := common.ValidatorAccept{ 980, 2.0,1, 2.0}
			accept5 := common.ValidatorAccept{ 980, 5.0,1, 5.0}
			err = dbval.DepositNote(123, &accept1)
			if err != nil{
				log.Error("Can't deposit notes: %s", err)
			}
			err = dbval.DepositNote(124, &accept2)
			if err != nil{
				log.Error("Can't deposit notes: %s", err)
			}
			err = dbval.DepositNote(125, &accept5)
			if err != nil{
				log.Error("Can't deposit notes: %s", err)
			}

			err := dbval.ReadNoteList(data)
			if err != nil{
				log.Error("Can't read notes: %s", err)
			}
			log.Info(data.String())
		}
		err = store.Close()
	}

	//err = linker.GetLinkerPorts(log)
	//
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
