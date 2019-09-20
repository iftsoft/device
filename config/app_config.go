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
	str := fmt.Sprintf("Duplex app config: \n%s\n%s\nDevice config: %v",
		cfg.Logger.String(), cfg.Duplex.String(), cfg.Device)
	return str
}

func GetDefaultAppConfig() *AppConfig {
	appCfg := &AppConfig{
		Logger: core.GetDefaultConfig(""),
		Duplex: duplex.GetDefaultClientConfig(),
		Device: GetDefaultDeviceConfig(),
	}
	return appCfg
}

func GetAppConfig(appPar *AppParams) (error, *AppConfig) {
	appCfg := GetDefaultAppConfig()
	err := core.ReadYamlFile(appPar.Config, appCfg)
	if err != nil {
		return err, nil
	} else {
		appPar.UpdateLoggerConfig(appCfg.Logger)
		appPar.UpdateClientConfig(appCfg.Duplex)
	}
	return nil, appCfg
}
