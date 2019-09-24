package main

import (
	"fmt"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/handler"
	"github.com/iftsoft/device/linker"
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

	err = GetLinkerPorts(log)

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

func GetLinkerPorts(out *core.LogAgent) error {
	out.Info("Serial ports")
	serList, err := linker.EnumerateSerialPorts()
	if err == nil {
		for i, ser := range serList {
			out.Info("   Port#%d - %s", i, ser)
		}
	}
	out.Info("HID / USB ports")
	hidList, err := linker.EnumerateHidUsbPorts()
	if err == nil {
		for i, hid := range hidList {
			out.Info("   Port#%d - %d:%d/%s", i, hid.VendorID, hid.ProductID, hid.Serial)
		}
	}
	return err
}
