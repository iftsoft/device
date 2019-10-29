package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func ReadYamlFile(name string, cfg interface{}) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckOrCreateFile(filename string) error {
	if FileExists(filename) {
		return nil
	}
	file, err := os.Create(filename)
	if err == nil {
		err = file.Close()
	}
	return err
}

