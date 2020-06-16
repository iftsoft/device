package handler

import (
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/duplex"
)

type HandlerManager struct {
	reflex ReflexSet
	router HandlerRouter
	proxy  HandlerProxy
	runner BinaryLauncher
}


func NewHandlerManager(config config.HandlerList) *HandlerManager {
	hm := &HandlerManager{}
	hm.reflex.initReflexSet()
	hm.router.initRouter(config)
	hm.proxy.initProxy()
	hm.runner.initLauncher(config)
	return hm
}


func (hm *HandlerManager) SetupDuplexServer(server duplex.ServerManager) {
	hm.proxy.setupProxy(server, &hm.router)
}


// Implementation of duplex.ClientManager
func (hm *HandlerManager) OnClientStarted(name string, info *duplex.GreetingInfo) {
	if name == "" {
		return
	}
	handler := hm.router.onClientStarted(name)
	if handler != nil {
		hm.reflex.attachReflexes(handler, &hm.proxy, info)
		handler.AttachProxy(&hm.proxy)
		handler.OnClientStarted(name)
	}
}

func (hm *HandlerManager) OnClientStopped(name string) {
	if name == "" {
		return
	}
	hm.router.onClientStopped(name)
}


func (hm *HandlerManager) RegisterReflexFactory(factory ReflexCreator) {
	hm.reflex.registerFactory(factory)
}

func (hm *HandlerManager) LaunchAllBinaries() {
	hm.runner.launchAllBinaries()
}

func (hm *HandlerManager) StopAllBinaries() {
	hm.runner.setQuitFlag()
	hm.router.terminateAll(&hm.proxy)
	hm.runner.waitAll()
}


func (hm *HandlerManager) Cleanup() {
	hm.router.cleanupRouter()
}


