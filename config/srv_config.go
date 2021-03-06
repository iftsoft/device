package config

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SrvConfig struct {
	Logger   *core.LogConfig		`yaml:"logger"`
	Duplex   *duplex.ServerConfig	`yaml:"duplex"`
	Handlers HandlerList			`yaml:"handlers"`
}

func (cfg *SrvConfig) String() string {
	str := fmt.Sprintf("Server app config: %s %s %s",
		cfg.Logger, cfg.Duplex, cfg.Handlers)
	return str
}

func GetDefaultSrvConfig() *SrvConfig {
	appCfg := &SrvConfig{
		Logger: core.GetDefaultConfig(""),
		Duplex: duplex.GetDefaultServerConfig(),
	}
	return appCfg
}

func GetSrvConfig(appPar *AppParams) (error, *SrvConfig) {
	appCfg := GetDefaultSrvConfig()
	err := core.ReadYamlFile(appPar.Config, appCfg)
	if err != nil {
		return err, nil
	} else {
		appPar.UpdateLoggerConfig(appCfg.Logger)
	}
	return nil, appCfg
}
