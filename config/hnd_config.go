package config

import "fmt"

type CommandConfig struct {
	DeviceName  string	`yaml:"device_name"`
	Enabled     bool    `yaml:"enabled"`
	BinaryFile  string	`yaml:"binary_file"`
	ConfigFile  string	`yaml:"config_file"`
	LoggerPath  string	`yaml:"logger_path"`
	Database    string	`yaml:"database"`
}
func (cfg *CommandConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tCommand config: " +
		"DeviceName = %s, Enabled = %t, BinaryFile = %s, ConfigFile = %s, LoggerPath = %s, Database = %s.",
		cfg.DeviceName, cfg.Enabled, cfg.BinaryFile, cfg.ConfigFile, cfg.LoggerPath, cfg.Database)
	return str
}

type ReflexConfig struct {
	ReflexName string `yaml:"reflex_name"`
	Enabled    bool   `yaml:"enabled"`
}
func (cfg *ReflexConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\t\t" +
		"ReflexName = %s, Enabled = %t.",
		cfg.ReflexName, cfg.Enabled)
	return str
}

type ReflexList []*ReflexConfig
func (cfg ReflexList) String() string {
	if cfg == nil { return "" }
	str := ""
	for _, plug := range cfg {
		str += plug.String()
	}
	return str
}

type HandlerConfig struct {
	Command  CommandConfig    `yaml:"command"`
	Config   ConfigOverwrite  `yaml:"config"`
	Reflexes ReflexList       `yaml:"reflexes"`
}
func (cfg *HandlerConfig) String() string {
	str := fmt.Sprintf("\n\tHandler config: %s %s %s",
		cfg.Command.String(), cfg.Config.String(), cfg.Reflexes)
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
