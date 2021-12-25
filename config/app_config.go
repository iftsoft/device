package config

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/duplex"
)

type AppConfig struct {
	Logger  *core.LogConfig      `yaml:"logger"`
	Duplex  *duplex.ClientConfig `yaml:"duplex"`
	Storage *dbase.StorageConfig `yaml:"storage"`
	Devices DeviceConfigList     `yaml:"devices"`
}

func (cfg *AppConfig) String() string {
	str := fmt.Sprintf("Client app config: %s %s %s %s",
		cfg.Logger, cfg.Duplex, cfg.Storage, cfg.Devices)
	return str
}

func (cfg *AppConfig) UpdateConfig(appPar *AppParams) {
	if cfg == nil {
		return
	}
	if cfg.Logger.LogFile == "" {
		cfg.Logger.LogFile = appPar.Name
	}
	if appPar.Logs != "" {
		cfg.Logger.LogPath = appPar.Logs
	}
	if appPar.DBase != "" {
		cfg.Storage.FileName = appPar.DBase
	}
}

func GetDefaultAppConfig() *AppConfig {
	appCfg := &AppConfig{
		Logger:  core.GetDefaultConfig(""),
		Duplex:  duplex.GetDefaultClientConfig(),
		Storage: dbase.GetDefaultStorageConfig(),
		Devices: GetDefaultDeviceConfigList(),
	}
	return appCfg
}

func GetAppConfig(appPar *AppParams) (error, *AppConfig) {
	appCfg := GetDefaultAppConfig()
	err := core.ReadYamlFile(appPar.Config, appCfg)
	if err != nil {
		return err, nil
	} else {
		appCfg.UpdateConfig(appPar)
	}
	return nil, appCfg
}
