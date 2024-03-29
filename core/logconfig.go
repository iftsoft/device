package core

import (
	"errors"
	"fmt"
)

type LogConfig struct {
	LogPath   string       `yaml:"log_path"`
	LogFile   string       `yaml:"log_file"`
	LogLevel  EnumLogLevel `yaml:"log_level"`
	ConsLevel EnumLogLevel `yaml:"cons_level"`
	MaxFiles  int          `yaml:"max_files"` // limit the number of log files under `logPath`
	MaxSize   int64        `yaml:"max_size"`  // limit size of a log file (KByte)
	Source    bool         `yaml:"source"`
}

func (cfg *LogConfig) PrintData() {
	fmt.Println("LogPath  ", cfg.LogPath)
	fmt.Println("LogFile  ", cfg.LogFile)
	fmt.Println("LogLevel ", cfg.LogLevel.String())
	fmt.Println("ConsLevel", cfg.ConsLevel.String())
	fmt.Println("MaxFiles ", cfg.MaxFiles)
	fmt.Println("MaxSize  ", cfg.MaxSize)
	fmt.Println("Source   ", cfg.Source)
}
func (cfg *LogConfig) String() string {
	if cfg == nil {
		return ""
	}
	str := fmt.Sprintf("\nLogging config: "+
		"LogPath = %s, LogFile = %s, LogLevel = %s, ConsLevel = %s, MaxFiles = %d, MaxSize = %d, Source = %v.",
		cfg.LogPath, cfg.LogFile, cfg.LogLevel, cfg.ConsLevel, cfg.MaxFiles, cfg.MaxSize, cfg.Source)
	return str
}

func GetDefaultConfig(name string) *LogConfig {
	cfg := LogConfig{
		LogPath:   "",
		LogFile:   name,
		LogLevel:  LogLevelTrace,
		ConsLevel: LogLevelDump,
		MaxFiles:  8,
		MaxSize:   1024,
		Source:    false,
	}
	return &cfg
}

func checkLogConfig(cfg *LogConfig) (err error) {
	if cfg == nil {
		return errors.New("logging: config is not set")
	}
	if cfg.LogFile == "" {
		return errors.New("logging: file name is not set")
	}
	if cfg.LogPath == "" {
		cfg.LogPath = "."
	}
	if cfg.LogLevel < LogLevelEmpty || cfg.LogLevel > LogLevelTrace {
		cfg.LogLevel = LogLevelInfo
	}
	if cfg.ConsLevel < LogLevelEmpty || cfg.ConsLevel > LogLevelTrace {
		cfg.ConsLevel = LogLevelError
	}
	if cfg.MaxFiles < 0 || cfg.MaxFiles >= 1024 {
		cfg.MaxFiles = 8
	}
	if cfg.MaxSize < 0 || cfg.MaxSize >= 128*1024 {
		cfg.MaxSize = 1024
	}
	return err
}
