package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type CfgLoadFilesYamlFolder struct {
	Path      string `yaml:"path"`
	Recursive bool   `yaml:"recursive"`
}

type CfgLoadFilesYaml struct {
	Folders         []*CfgLoadFilesYamlFolder `yaml:"folders"`
	AllowedSuffixes []string                  `yaml:"allowed_suffixes"`
}

type LoadFilesYamlOptions struct {
	Logger       *log.Entry
	DataProvider func() (data interface{})
	Config       *CfgLoadFilesYaml
}

func LoadFilesYaml(opts *LoadFilesYamlOptions) (datas []interface{}, errs []error) {
	logger := opts.Logger.WithFields(log.Fields{"Function": "Load", "SourceFolders": len(opts.Config.Folders)})
	logger.Infof("Loading files")
	for _, cfgFolder := range opts.Config.Folders {
		logger := logger.WithFields(log.Fields{"FolderPath": cfgFolder.Path})
		logger.Infof("Loading dashboard config from folder")
		data, err := loadFilesYamlFromFolder(cfgFolder, opts)
		if len(err) != 0 {
			errs = append(errs, err...)
		}
		if len(data) > 0 {
			datas = append(datas, data...)
		}
	}
	return datas, errs
}

func loadFilesYamlHasAllowedSuffix(name string, opts *LoadFilesYamlOptions) (hasAllowedSuffix bool) {
	allowedSuffixes := opts.Config.AllowedSuffixes
	if len(allowedSuffixes) == 0 {
		allowedSuffixes = []string{".yaml", ".yml"}
	}
	logger := opts.Logger.WithFields(log.Fields{"Function": "loadFilesYamlHasAllowedSuffix", "Name": name, "AllowedSuffixes": allowedSuffixes})
	for _, allowedSuffix := range allowedSuffixes {
		if strings.HasSuffix(name, allowedSuffix) {
			return true
		}
	}
	logger.Warnf("This file does not have an allowed yaml file name suffix")
	return false
}

func loadFilesYamlFromFolder(cfgFolder *CfgLoadFilesYamlFolder, opts *LoadFilesYamlOptions) (datas []interface{}, errs []error) {
	logger := opts.Logger.WithFields(log.Fields{"Function": "loadFromFolder", "ConfigFolderPath": cfgFolder.Path})
	if cfgFolder.Recursive {
		return loadFilesYamlFromFolderRecursive(cfgFolder, opts)
	}
	entries, err := os.ReadDir(cfgFolder.Path)
	if err != nil {
		return nil, []error{fmt.Errorf("when reading the directory '%s': %w", cfgFolder.Path, err)}
	}
	for idx := range entries {
		entry := entries[idx]
		logger := logger.WithFields(log.Fields{"FileName": entry.Name()})
		if entry.IsDir() {
			continue
		}
		if !loadFilesYamlHasAllowedSuffix(entry.Name(), opts) {
			continue
		}
		if data, err := loadFilesYamlFromFile(filepath.Join(cfgFolder.Path, entry.Name()), opts); err != nil {
			logger.Errorf("when loading the dashboard from file: %s", err)
			errs = append(errs, err)
		} else {
			datas = append(datas, data)
		}
	}
	return datas, errs
}

func loadFilesYamlFromFolderRecursive(cfgFolder *CfgLoadFilesYamlFolder, opts *LoadFilesYamlOptions) (datas []interface{}, errs []error) {
	if err := filepath.WalkDir(cfgFolder.Path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if !loadFilesYamlHasAllowedSuffix(d.Name(), opts) {
			return nil
		}

		if data, err := loadFilesYamlFromFile(path, opts); err != nil {
			errs = append(errs, fmt.Errorf("when loading data: %w", err))
		} else {
			datas = append(datas, data)
		}
		return nil
	}); err != nil {
		errs = append(errs, fmt.Errorf("when walking the dir starting at '%s': %w", cfgFolder.Path, err))
	}
	return datas, errs
}

func loadFilesYamlFromFile(path string, opts *LoadFilesYamlOptions) (data interface{}, err error) {
	logger := opts.Logger.WithFields(log.Fields{"Function": "loadFromFolder", "Path": path})
	logger.Infof("Loading data from file")
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("when reading the file '%s': %w", path, err)
	}
	data = opts.DataProvider()
	if err := yaml.Unmarshal(yamlFile, data); err != nil {
		return nil, fmt.Errorf("when unmarshaling the file '%s': %w", path, err)
	}
	logger.Infof("Loaded")
	return data, nil
}
