package postgres

import (
	"fmt"
	"strings"
	"time"
)

type Postgres struct {
	Name        string        `json:"name" yaml:"name"`
	Host        string        `json:"host" yaml:"host"`
	Port        int           `json:"port" yaml:"port"`
	Username    string        `json:"username" yaml:"username"`
	Password    string        `json:"password" yaml:"password"`
	Database    string        `json:"database" yaml:"database"`
	SSLMode     string        `json:"sslMode" yaml:"sslMode"`
	SSLCertPath string        `json:"sslCertPath" yaml:"sslCertPath"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`

	TableAlerts string `json:"tableAlerts" yaml:"tableAlerts"`
	TableKV     string `json:"tableKV" yaml:"tableKV"`
}

func (cfg *Postgres) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return fmt.Errorf("host must be defined")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	if strings.TrimSpace(cfg.TableAlerts) == "" {
		return fmt.Errorf("empty tableAlerts")
	}
	if strings.TrimSpace(cfg.TableKV) == "" {
		return fmt.Errorf("empty tableKV")
	}

	return nil
}
