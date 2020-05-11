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

	//err = core.CheckOrCreateFile(appCfg.Storage.FileName)
	//if err != nil{
	//	log.Error("Can't check database file: %s", err)
	//} else {
	//	store := dbase.GetNewDBaseStore(appCfg.Storage)
	//	err = store.Open()
	//	if err != nil{
	//		log.Error("Can't open database file: %s", err)
	//	} else {
	//		dbval := dbvalid.NewDBaseValidator(store, appCfg.Duplex.DevName)
	//		err = dbval.CreateAllTables()
	//
	//		data := &common.ValidatorBatch{}
	//		err = dbval.CloseBatch(data)
	//
	//		defNoteList := common.ValidNoteList {
	//			{appCfg.Duplex.DevName, 980, 0, 1.0, 0.0, },
	//			{appCfg.Duplex.DevName, 980, 0, 2.0, 0.0, },
	//			{appCfg.Duplex.DevName, 980, 0, 5.0, 0.0, },
	//		}
	//		err = dbval.InitNoteList(defNoteList)
	//
	//		accept1 := common.ValidatorAccept{ 980, 1.0,1, 1.0}
	//		accept2 := common.ValidatorAccept{ 980, 2.0,1, 2.0}
	//		accept5 := common.ValidatorAccept{ 980, 5.0,1, 5.0}
	//		err = dbval.DepositNote(123, &accept1)
	//		err = dbval.DepositNote(124, &accept2)
	//		err = dbval.DepositNote(125, &accept5)
	//
	//		err = dbval.ReadNoteList(data)
	//	}
	//	err = store.Close()
	//}

	err = linker.GetLinkerPorts(log)

	err = CheckValidatorDatabase(appCfg.Storage, appCfg.Duplex.DevName)
	if err == nil{
		dev := driver.NewSystemDevice(appCfg)
		drv := validator.NewValidatorDriver()
		err = dev.InitDevice(drv)
		if err == nil {
			dev.StartDevice()

			core.WaitForSignal(log)

			dev.StopDevice()
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
