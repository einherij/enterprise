package db

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"database"`
}

func NewPostgresClient(cfg PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	pgDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			logrus.WithField("component", "gorm"),
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      false,
				LogLevel:      logger.Warn,
			}),
	})
	if err != nil {
		return nil, fmt.Errorf(
			"cannot connect to the postgresql database %q at %s:%s: %+v",
			cfg.DBName, cfg.Host, cfg.Port, err,
		)
	}
	logrus.Infof("connected to the database %q at %s:%s", cfg.DBName, cfg.Host, cfg.Port)
	return pgDB, nil
}
