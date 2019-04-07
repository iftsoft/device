package main

import (
	"fmt"
	"github.com/iftsoft/core"
	"github.com/iftsoft/device/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppConfig struct {
	Logger core.LogConfig      `yaml:"logger"`
	Device config.DeviceConfig `yaml:"device"`
}

func main() {
	fmt.Println("-------BEGIN------------")

	config := &AppConfig{}
	err := core.ReadYamlFile("config.yml", config)
	if err != nil {
		fmt.Println(err)
	} else {
		core.StartFileLogger(&config.Logger)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start application")
	log.Info("Config logger: %+v", config.Logger)
	log.Info("Config device: %+v", config.Device)

	WaitForSignal(log)

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
