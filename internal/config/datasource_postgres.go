package config

import (
	"fmt"
	"strings"
)

type DataSourcePostgres struct {
	Name        string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Database    string `json:"database" yaml:"database"`
	SSLMode     string `json:"sslMode" yaml:"sslMode"`
	SSLCertPath string `json:"sslCertPath" yaml:"sslCertPath"`
}

func (cfg DataSourcePostgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if cfg.Host == "" {
		return fmt.Errorf("host must be defined")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}

	return nil
}
