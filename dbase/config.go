package dbase

import "fmt"

type StorageConfig struct {
	FileName   string `yaml:"file_name"`
}
func (cfg *StorageConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\nStorage config: " +
		"FileName = %s.",
		cfg.FileName)
	return str
}
func GetDefaultStorageConfig() *StorageConfig {
	cfg := &StorageConfig{
		FileName:    "",
	}
	return cfg
}


