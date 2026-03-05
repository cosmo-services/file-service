package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DirectoryConfig struct {
	Access string `yaml:"access"`
	Path   string `yaml:"path"`
}

type FileStorageConfig struct {
	Directories map[string]DirectoryConfig `yaml:"directory"`
}

func NewFileStorageConfig() (*FileStorageConfig, error) {
	cfg, err := LoadFileStorageConfig("./storage.yaml")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadFileStorageConfig(path string) (*FileStorageConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config FileStorageConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func (c *FileStorageConfig) GetDirectoryInfo(logicalDir string) (*DirectoryConfig, error) {
	dir, ok := c.Directories[logicalDir]
	if !ok {
		return nil, fmt.Errorf("directory not configured: %s", logicalDir)
	}
	return &dir, nil
}

func (c *FileStorageConfig) IsPublic(logicalDir string) (bool, error) {
	dir, err := c.GetDirectoryInfo(logicalDir)
	if err != nil {
		return false, err
	}
	return dir.Access == "public", nil
}

func (c *FileStorageConfig) GetPhysicalPath(logicalDir string) (string, error) {
	dir, err := c.GetDirectoryInfo(logicalDir)
	if err != nil {
		return "", err
	}
	return dir.Path, nil
}

func (c *FileStorageConfig) GetPublicDirectoryNames() []string {
	var names []string

	for name, dir := range c.Directories {
		if dir.Access == "public" {
			names = append(names, name)
		}
	}

	return names
}
