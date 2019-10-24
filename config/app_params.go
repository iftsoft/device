package config

import (
	"flag"
	"fmt"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"os"
	"path/filepath"
	"strings"
)

type AppParams struct {
	Home   string   // Working folder for application
	Name   string   // Name of application
	Config string   // Application config file
	DBase  string   // Path to database file
	Logs   string   // Path to log files folder
	Args   []string // Rest of application params
}

func GetAppParams() *AppParams {
	path, name := filepath.Split(os.Args[0])
	full, err := filepath.Abs(path)
	if err == nil {
		path = full
	}
	name = strings.TrimSuffix(name, filepath.Ext(name))
	conf := path + string(os.PathSeparator) + name + ".yml"
	logs := path + string(os.PathSeparator) + "logs"
	base := path + string(os.PathSeparator) + name + ".db"
	appPar := AppParams{}
	flag.StringVar(&appPar.Home, "home", path, "Working folder for application")
	flag.StringVar(&appPar.Name, "name", name, "Name of application")
	flag.StringVar(&appPar.Config, "cfg", conf, "Application config file")
	flag.StringVar(&appPar.DBase, "base", base, "Path to database file")
	flag.StringVar(&appPar.Logs, "logs", logs, "Path to log files folder")
	// Parse command line
	flag.Parse()
	// Get rest of params
	appPar.Args = flag.Args()
	//	appPar.PrintData()
	return &appPar
}

func (par *AppParams) PrintData() {
	fmt.Println("Home dir ", par.Home)
	fmt.Println("App name ", par.Name)
	fmt.Println("Config   ", par.Config)
	fmt.Println("Logs dir ", par.Logs)
	fmt.Println("Database ", par.DBase)
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
	if cfg.LogFile == "" {
		cfg.LogFile = par.Name
	}
	if cfg.LogPath == "" {
		cfg.LogPath = par.Logs
	}
}

func (par *AppParams) UpdateClientConfig(cfg *duplex.ClientConfig) {
	if cfg == nil {
		return
	}
	if cfg.DevName == "" {
		cfg.DevName = par.Name
	}
}

func (par *AppParams) UpdateStorageConfig(cfg *StorageConfig) {
	if cfg == nil {
		return
	}
	if cfg.FileName == "" {
		cfg.FileName = par.DBase
	}
}
