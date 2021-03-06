package sql

import (
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type PostgresAlert struct {
	db       *sqlx.DB
	tableCfg tables.TableAlerts
	timeout  time.Duration
	logger   *zap.Logger
}
