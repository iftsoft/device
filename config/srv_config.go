package config

import (
	"flag"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"os"
	"path/filepath"
	"strings"
)

type SrvConfig struct {
	Logger core.LogConfig      `yaml:"logger"`
	Server duplex.ServerConfig `yaml:"server"`
}

func GetDefaultSrvConfig() *SrvConfig {
	appCfg := &SrvConfig{
		Logger: core.LogConfig{
			LogPath:   "",
			LogFile:   "",
			LogLevel:  core.LogLevelTrace,
			ConsLevel: core.LogLevelTrace,
			MaxFiles:  8,
			MaxSize:   1024,
		},
		Server: duplex.ServerConfig{
			Port: duplex.DuplexPort,
		},
	}
	return appCfg
}

type SrvParams struct {
	Home   string   // Working folder for application
	Name   string   // Name of application
	Config string   // Application config file
	Logs   string   // Path to log files folder
	Args   []string // Rest of application params
}

func GetSrvParams() *SrvParams {
	path, name := filepath.Split(os.Args[0])
	name = strings.TrimSuffix(name, filepath.Ext(name))
	appPar := SrvParams{}
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

func GetSrvConfig(appPar *SrvParams) (error, *SrvConfig) {
	appCfg := GetDefaultSrvConfig()
	err := core.ReadYamlFile(appPar.Config, appCfg)
	if err != nil {
		return err, nil
	} else {
		UpdateSrvConfig(appCfg, appPar)
	}
	return nil, appCfg
}

func UpdateSrvConfig(appCfg *SrvConfig, appPar *SrvParams) {
	if appCfg.Logger.LogFile == "" {
		appCfg.Logger.LogFile = appPar.Name
	}
	if appCfg.Logger.LogPath == "" {
		appCfg.Logger.LogPath = appPar.Logs
	}
}
