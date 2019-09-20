package config

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SrvConfig struct {
	Logger *core.LogConfig      `yaml:"logger"`
	Duplex *duplex.ServerConfig `yaml:"duplex"`
}

func (cfg *SrvConfig) String() string {
	str := fmt.Sprintf("Duplex app config: \n%s\n%s",
		cfg.Logger.String(), cfg.Duplex.String())
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
