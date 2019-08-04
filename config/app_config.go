package config

import (
	"flag"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"os"
	"path/filepath"
	"strings"
)

type AppConfig struct {
	Logger core.LogConfig            `yaml:"logger"`
	Client duplex.DuplexClientConfig `yaml:"client"`
	Device DeviceConfig              `yaml:"device"`
}

func GetDefaultAppConfig() *AppConfig {
	appCfg := &AppConfig{
		Logger: core.LogConfig{
			LogPath:   "",
			LogFile:   "",
			LogLevel:  core.LogLevelTrace,
			ConsLevel: core.LogLevelTrace,
			MaxFiles:  8,
			MaxSize:   1024,
		},
		Client: duplex.DuplexClientConfig{
			Port:    duplex.DuplexPort,
			DevName: "TestDevice",
		},
		Device: DeviceConfig{
			Common:    CommonConfig{},
			Serial:    SerialConfig{},
			Printer:   PrinterConfig{},
			Reader:    ReaderConfig{},
			Validator: ValidatorConfig{},
			Dispenser: DispenserConfig{},
			Vendor:    VendorConfig{},
		},
	}
	return appCfg
}

type AppParams struct {
	Home   string   // Working folder for application
	Name   string   // Name of application
	Config string   // Application config file
	Logs   string   // Path to log files folder
	Args   []string // Rest of application params
}

func GetAppParams() *AppParams {
	path, name := filepath.Split(os.Args[0])
	name = strings.TrimSuffix(name, filepath.Ext(name))
	appPar := AppParams{}
	flag.StringVar(&appPar.Home, "home", path, "Working folder for application")
	flag.StringVar(&appPar.Name, "name", name, "Name of application")
	flag.StringVar(&appPar.Config, "cfg", name+".yml", "Application config file")
	flag.StringVar(&appPar.Logs, "logs", "logs", "Path to log files folder")
	// Parse command line
	flag.Parse()
	// Get rest of params
	appPar.Args = flag.Args()
	return &appPar
}

func GetAppConfig(appPar *AppParams) (error, *AppConfig) {
	appCfg := GetDefaultAppConfig()
	err := core.ReadYamlFile(appPar.Config, appCfg)
	if err != nil {
		return err, nil
	} else {
		UpdateAppConfig(appCfg, appPar)
	}
	return nil, appCfg
}

func UpdateAppConfig(appCfg *AppConfig, appPar *AppParams) {
	if appCfg.Logger.LogFile == "" {
		appCfg.Logger.LogFile = appPar.Name
	}
	if appCfg.Logger.LogPath == "" {
		appCfg.Logger.LogPath = appPar.Logs
	}
}
