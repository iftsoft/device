package config

import (
	"flag"
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/duplex"
	"os"
	"path/filepath"
	"strings"
)

type AppParams struct {
	Name     string   // Name of application
	Home     string   // Working folder for application
	Config   string   // Application config file
	DBase    string   // Path to database file
	Logs     string   // Path to log files folder
	Args     []string // Rest of application params
	defName  string   // Name of application
	defHome  string   // Working folder for application
	defConf  string   // Application config file
	defBase  string   // Path to database file
	defLogs  string   // Path to log files folder
}

func GetAppParams() *AppParams {
	appPar := AppParams{}
	path, name := filepath.Split(os.Args[0])
	full, err := filepath.Abs(path)
	if err == nil {
		path = full
	}
	appPar.defHome = path
	appPar.defName = strings.TrimSuffix(name, filepath.Ext(name))
	appPar.defConf = path + string(os.PathSeparator) + appPar.defName + ".yml"
	appPar.defLogs = path + string(os.PathSeparator) + "logs"
	appPar.defBase = path + string(os.PathSeparator) + appPar.defName + ".db"
	var parName, parHome, parCfg, parBase, parLogs string
	flag.StringVar(&parName, "name", "", "Name of application")
	flag.StringVar(&parHome, "home", "", "Working folder for application")
	flag.StringVar(&parCfg,  "cfg",  "", "Application config file")
	flag.StringVar(&parBase, "base", "", "Path to database file")
	flag.StringVar(&parLogs, "logs", "", "Path to log files folder")
	// Parse command line
	flag.Parse()
	// Get rest of params
	appPar.Name   = strings.Trim(parName, "\"")
	appPar.Home   = strings.Trim(parHome, "\"")
	appPar.Config = strings.Trim(parCfg,  "\"")
	appPar.DBase  = strings.Trim(parBase, "\"")
	appPar.Logs   = strings.Trim(parLogs, "\"")
	appPar.Args   = flag.Args()
	if appPar.Config == "" {
		appPar.Config = appPar.defConf
	}

//	appPar.PrintData()
	return &appPar
}

func (par *AppParams) PrintData() {
	fmt.Println("App name ", par.Name)
	fmt.Println("Home dir ", par.Home)
	fmt.Println("Config   ", par.Config)
	fmt.Println("Database ", par.DBase)
	fmt.Println("Logs dir ", par.Logs)
	fmt.Println("Args     ", par.Args)
}

func (par *AppParams) String() string {
	str := fmt.Sprintf("App params: "+
		"Home = %s, Name = %s, Config = %s, DBase = %s, Logs = %s, Args = %v.",
		par.Home, par.Name, par.Config, par.DBase, par.Logs, par.Args)
	return str
}

func (par *AppParams) UpdateLoggerConfig(cfg *core.LogConfig) {
	if cfg == nil {
		return
	}
	if par.Name != "" {
		cfg.LogFile = par.Name
	}
	if cfg.LogFile == "" {
		cfg.LogFile = par.defName
	}

	if par.Logs != "" {
		cfg.LogPath = par.Logs
	}
	if cfg.LogPath == "" {
		cfg.LogPath = par.defLogs
	}
}

func (par *AppParams) UpdateClientConfig(cfg *duplex.ClientConfig) {
	if cfg == nil {
		return
	}
	if par.Name != "" {
		cfg.DevName = par.Name
	}
	if cfg.DevName == "" {
		cfg.DevName = par.defName
	}
}

func (par *AppParams) UpdateStorageConfig(cfg *dbase.StorageConfig) {
	if cfg == nil {
		return
	}
	if par.DBase != "" {
		cfg.FileName = par.DBase
	}
	if cfg.FileName == "" {
		cfg.FileName = par.defBase
	}
}
