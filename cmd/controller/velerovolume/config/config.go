package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type VeleroVolumeCfg struct {
	IncludeNamespaces  string `yaml:"includeNamespaces,omitempty"`
	ExcludeNamespaces  string `yaml:"excludeNamespaces,omitempty"`
	IncludeVolumeTypes string `yaml:"includeVolumeTypes,omitempty"`
	ExcludeVolumeTypes string `yaml:"excludeVolumeTypes,omitempty"`
}

type ClusterServerCfg struct {
	// The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.
	MasterURL string `yaml:"masterURL,omitempty"`
	// Path to a kubeconfig. Only required if out-of-cluster.
	KubeConfig string `yaml:"kubeConfig,omitempty"`
	// LeaseLock namespace
	LeaseLockNamespace string `yaml:"leaseLockNamespace,omitempty"`
	// LeaseLock name
	LeaseLockName string `yaml:"leaseLockName,omitempty"`
}

type Config struct {
	ClusterServerCfg *ClusterServerCfg `yaml:"clusterServerCfg,omitempty"`
	VeleroVolumeCfg  *VeleroVolumeCfg  `yaml:"veleroVolumeCfg,omitempty"`
}

// validate the configuration
func (c *Config) validate() error {
	if c.VeleroVolumeCfg.IncludeNamespaces != "" && c.VeleroVolumeCfg.ExcludeNamespaces != "" ||
		c.VeleroVolumeCfg.IncludeVolumeTypes != "" && c.VeleroVolumeCfg.ExcludeVolumeTypes != "" {
		return fmt.Errorf("Invalid velero volume resources configurations, please check ...")
	}
	// TODO: other configuration validate ...
	return nil
}

// LoadConfig parses configuration file and returns
// an initialized Settings object and an error object if any. For instance if it
// cannot find the configuration file it will set the returned error appropriately.
func LoadConfig(path string) (*Config, error) {
	c := &Config{}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read configuration file: %s,error: %s", path, err)
	}
	if err = yaml.Unmarshal(contents, c); err != nil {
		return nil, fmt.Errorf("Failed to parse configuration,error: %s", err)
	}
	if err = c.validate(); err != nil {
		return nil, fmt.Errorf("Invalid configuration,error: %s", err)
	}
	return c, nil
}
