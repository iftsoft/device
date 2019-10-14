package config

import "fmt"

type CommandConfig struct {
	DeviceName  string	`yaml:"device_name"`
	BinaryFile  string	`yaml:"binary_file"`
	ConfigFile  string	`yaml:"config_file"`
	LoggerPath  string	`yaml:"logger_path"`
	Database    string	`yaml:"database"`
	AutoLoad	bool	`yaml:"auto_load"`
}
func (cfg *CommandConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tCommand config: " +
		"DeviceName = %s, BinaryFile = %s, ConfigFile = %s, LoggerPath = %s, Database = %s, AutoLoad = %t.",
		cfg.DeviceName, cfg.BinaryFile, cfg.ConfigFile, cfg.LoggerPath, cfg.Database, cfg.AutoLoad)
	return str
}

type PluginConfig struct {
	PluginName  string	`yaml:"plugin_name"`
	Enabled     bool	`yaml:"enabled"`
}
func (cfg *PluginConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\t\t" +
		"PluginName = %s, Enabled = %t.",
		cfg.PluginName, cfg.Enabled)
	return str
}

type PluginList []*PluginConfig
func (cfg PluginList) String() string {
	if cfg == nil { return "" }
	str := ""
	for _, plug := range cfg {
		str += plug.String()
	}
	return str
}

type HandlerConfig struct {
	Command  CommandConfig	`yaml:"command"`
	Plugins  PluginList		`yaml:"plugins"`
}
func (cfg *HandlerConfig) String() string {
	str := fmt.Sprintf("\n\tHandler config: %s %s",
		cfg.Command.String(), cfg.Plugins)
	return str
}

type HandlerList []*HandlerConfig
func (cfg HandlerList) String() string {
	if cfg == nil { return "" }
	str := "\nHandlers:"
	for _, hnd := range cfg {
		str += hnd.String()
	}
	return str
}
