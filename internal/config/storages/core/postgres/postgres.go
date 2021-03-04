package postgres

import (
	"fmt"
	"strings"
)

type Tables struct {
	Alerts string `json:"alerts" yaml:"alerts" hcl:"alerts"`
	KV     string `json:"kv" yaml:"kv" hcl:"kv"`
}

type Postgres struct {
	Name        string `json:"name" yaml:"name" hcl:"name,label"`
	Host        string `json:"host" yaml:"host" hcl:"host"`
	Port        int    `json:"port" yaml:"port" hcl:"port"`
	Username    string `json:"username" yaml:"username" hcl:"username"`
	Password    string `json:"password" yaml:"password" hcl:"password"`
	Database    string `json:"database" yaml:"database" hcl:"database"`
	SSLMode     string `json:"sslMode" yaml:"sslMode" hcl:"sslMode,optional"`
	SSLCertPath string `json:"sslCertPath" yaml:"sslCertPath" hcl:"sslCertPath,optional"`
	Timeout     int    `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`

	Tables Tables `json:"tables" yaml:"tables" hcl:"tables,block"`
}

func (cfg Postgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return fmt.Errorf("host must be defined")
	}
	if strings.TrimSpace(cfg.Tables.Alerts) == "" {
		return fmt.Errorf("table Alerts must be not empty")
	}
	if strings.TrimSpace(cfg.Tables.KV) == "" {
		return fmt.Errorf("table KV must be not empty")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
