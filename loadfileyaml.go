package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// LoadFileYaml will load a Yaml file and unmarshal it into the provided interface
func LoadFileYaml(logger *log.Entry, filename string, data interface{}) error {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading the file '%s': %w", filename, err)
	}
	err = yaml.Unmarshal(fileContent, data)
	if err != nil {
		return fmt.Errorf("unmarshal the file '%s': %w", filename, err)
	}
	return nil
}
