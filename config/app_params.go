package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AppParams struct {
	Name   string   // Name of application
	Config string   // Application config file
	DBase  string   // Path to database file
	Logs   string   // Path to log files folder
	Args   []string // Rest of application params
	// Parameters by default
	//defName string // Name of application by default
	//defConf string // Application config file by default
	//defBase string // Path to database file by default
	//defLogs string // Path to log files folder by default
}

func GetAppParams() *AppParams {
	appPar := AppParams{}
	path, name := filepath.Split(os.Args[0])
	full, err := filepath.Abs(path)
	if err == nil {
		path = full
	}
	defName := strings.TrimSuffix(name, filepath.Ext(name))
	defConf := path + string(os.PathSeparator) + defName + ".yml"
	defLogs := path + string(os.PathSeparator) + "logs"
	defBase := path + string(os.PathSeparator) + defName + ".db"
	var parName, parCfg, parBase, parLogs string
	flag.StringVar(&parName, "name", "", "Name of application")
	flag.StringVar(&parCfg, "cfg", "", "Application config file")
	flag.StringVar(&parBase, "base", "", "Path to database file")
	flag.StringVar(&parLogs, "logs", "", "Path to log files folder")
	// Parse command line
	flag.Parse()
	// Get rest of params
	appPar.Name = strings.Trim(parName, "\"")
	appPar.Config = strings.Trim(parCfg, "\"")
	appPar.DBase = strings.Trim(parBase, "\"")
	appPar.Logs = strings.Trim(parLogs, "\"")
	appPar.Args = flag.Args()
	if appPar.Name == "" {
		appPar.Name = defName
	}
	if appPar.Config == "" {
		appPar.Config = defConf
	}
	if appPar.DBase == "" {
		appPar.DBase = defBase
	}
	if appPar.Logs == "" {
		appPar.Logs = defLogs
	}
	//	appPar.PrintData()
	return &appPar
}

func (par *AppParams) PrintData() {
	fmt.Println("App name ", par.Name)
	fmt.Println("Config   ", par.Config)
	fmt.Println("Database ", par.DBase)
	fmt.Println("Logs dir ", par.Logs)
	fmt.Println("Add args ", par.Args)
}

func (par *AppParams) String() string {
	str := fmt.Sprintf("App params: "+
		"Name = %s, Config = %s, DBase = %s, Logs = %s, Args = %v.",
		par.Name, par.Config, par.DBase, par.Logs, par.Args)
	return str
}
