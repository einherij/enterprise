package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/sirupsen/logrus"
)

const (
	defaultOperationTimeout = 10 * time.Second
)

type ClickhouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Debug    bool   `mapstructure:"debug"`
}

func NewClickHouseClient(cfg ClickhouseConfig) (*sql.DB, error) {
	cxDSN := fmt.Sprintf(
		"tcp://%s:%s?username=%s&password=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
	)
	if cfg.Debug {
		cxDSN = fmt.Sprintf("%s&debug=true", cxDSN)
	}
	cxDB, err := sql.Open("clickhouse", cxDSN)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot connect to the clickhouse database at %s:%s: %w",
			cfg.Host, cfg.Port, err,
		)
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultOperationTimeout)
	defer cancel()
	if err = cxDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf(
			"clickhouse database is unavailable at %s:%s: %w",
			cfg.Host, cfg.Port, err,
		)
	}
	logrus.Infof("connected to the clickhouse database at %s:%s", cfg.Host, cfg.Port)
	return cxDB, nil
}
