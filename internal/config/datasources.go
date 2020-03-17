package config

type DataSources struct {
	Clickhouse []DataSourceClickhouse `json:"clickhouse" yaml:"clickhouse"`
	Prometheus []DataSourcePrometheus `json:"prometheus" yaml:"prometheus"`
	Postgres   []DataSourcePostgres   `json:"postgres" yaml:"postgres"`
	MySQL      []DataSourceMysql      `json:"mysql" yaml:"mysql"`
	Loki       []DataSourceLoki       `json:"loki" yaml:"loki"`
}

func (cfg DataSources) Validate() error {
	for _, c := range cfg.Clickhouse {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for _, c := range cfg.Prometheus {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for _, c := range cfg.Postgres {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for _, c := range cfg.MySQL {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for _, c := range cfg.Loki {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}
