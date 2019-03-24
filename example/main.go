package main

import (
	"fmt"
	"github.com/iftsoft/core/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	logCfg := log.GetDefaultConfig("example")
	log.StartFileLogger(logCfg)
	out := log.GetLogAgent(log.LogLevelTrace, "APP")
	out.Info("Start application")

	WaitForSignal(out)

	out.Info("Stop application")
	time.Sleep(time.Second)
	log.StopFileLogger()
	fmt.Println("-------END------------")
}

func WaitForSignal(out *log.LogAgent) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	out.Info("Got signal: %v, exiting.", s)
}
