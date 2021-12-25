package config

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SrvConfig struct {
	Logger   *core.LogConfig      `yaml:"logger"`
	Duplex   *duplex.ServerConfig `yaml:"duplex"`
	Handlers HandlerList          `yaml:"handlers"`
}

func (cfg *SrvConfig) String() string {
	str := fmt.Sprintf("Server app config: %s %s %s",
		cfg.Logger, cfg.Duplex, cfg.Handlers)
	return str
}

func (cfg *SrvConfig) UpdateConfig(appPar *AppParams) {
	if cfg == nil {
		return
	}
	if cfg.Logger.LogFile == "" {
		cfg.Logger.LogFile = appPar.Name
	}
	if appPar.Logs != "" {
		cfg.Logger.LogPath = appPar.Logs
	}
}

func GetDefaultSrvConfig() *SrvConfig {
	appCfg := &SrvConfig{
		Logger: core.GetDefaultConfig(""),
		Duplex: duplex.GetDefaultServerConfig(),
	}
	return appCfg
}

func GetSrvConfig(appPar *AppParams) (error, *SrvConfig) {
	srvCfg := GetDefaultSrvConfig()
	err := core.ReadYamlFile(appPar.Config, srvCfg)
	if err != nil {
		return err, nil
	} else {
		srvCfg.UpdateConfig(appPar)
	}
	return nil, srvCfg
}
