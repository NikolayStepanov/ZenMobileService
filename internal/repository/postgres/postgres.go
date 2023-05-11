package postgres

import (
	"ZenMobileService/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

const (
	driver = "postgres"
)

func NewPostgresDB(ctx context.Context, cfg *config.PostgresConfig) (*pgx.Conn, error) {
	connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		driver, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Error("Unable to parse config")
		log.Error(err)
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		log.Info("Unable to connect to database")
		log.Error(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Error(err)
		log.Infoln("Postgres is not connect")
		return nil, err
	} else {
		log.Infoln("Postgres is connected")
	}

	return conn, nil
}
