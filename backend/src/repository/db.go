package repository

import (
	"backend/src/core"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var sugar = core.Sugar
var ConObjOfDB *pgxpool.Pool

func DbPoolMain() {
	// Create database connection pool
	connPool, err := pgxpool.NewWithConfig(context.Background(), DbConnConfig())
	if err != nil {
		sugar.Errorf("Error while creating connection to the database!! %v", err.Error())
		return // Exit if there's an error
	}

	sugar.Info("Connected to the database!!")
	ConObjOfDB = connPool // Store the connection obj as Private
}

func DbConnConfig() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Database connection string
	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?connect_timeout=10",
		core.Config.Database.USER_NAME,
		core.Config.Database.PASSWORD,
		core.Config.Database.HOST,
		core.Config.Database.PORT,
		core.Config.Database.DATABASE)

	// Parse and create the config
	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		sugar.Fatalf("Failed to create a config, error: %v", err)
	}

	// Set pool configurations
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}
