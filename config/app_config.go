package config

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type AppConfig struct {
	Logger *core.LogConfig      `yaml:"logger"`
	Duplex *duplex.ClientConfig `yaml:"duplex"`
	Device *DeviceConfig        `yaml:"device"`
}

func (cfg *AppConfig) String() string {
	str := fmt.Sprintf("Client app config: %s %s %s",
		cfg.Logger, cfg.Duplex, cfg.Device)
	return str
}

func GetDefaultAppConfig(devCfg *DeviceConfig) *AppConfig {
	appCfg := &AppConfig{
		Logger: core.GetDefaultConfig(""),
		Duplex: duplex.GetDefaultClientConfig(),
		Device: devCfg,
	}
	return appCfg
}

func GetAppConfig(appPar *AppParams, devCfg *DeviceConfig) (error, *AppConfig) {
	appCfg := GetDefaultAppConfig(devCfg)
	err := core.ReadYamlFile(appPar.Config, appCfg)
	if err != nil {
		return err, nil
	} else {
		appPar.UpdateLoggerConfig(appCfg.Logger)
		appPar.UpdateClientConfig(appCfg.Duplex)
	}
	return nil, appCfg
}
