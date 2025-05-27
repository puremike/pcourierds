package main

import (
	"time"

	_ "github.com/lib/pq"
	"github.com/puremike/pcourierds/internal/db"
	"github.com/puremike/pcourierds/internal/env"
	"github.com/puremike/pcourierds/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config *config
	logger *zap.SugaredLogger
	store  *store.Storage
}

type config struct {
	port     string
	env      string
	dbconfig dbconfig
}

type dbconfig struct {
	db_url           string
	maxIdleConns     int
	maxOpenConns     int
	connsMaxIdleTime time.Duration
}

const apiVersion = "1.1.0"

// @title			Courier Delivery System API
// @version		1.1.0
// @description	This is an API for a Courier Delivery System
//
// @contact.name	Puremike
// @contact.url	http://github.com/puremike
// @contact.email	digitalmarketfy@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath		/api/v1
func main() {
	cfg := &config{
		port: env.GetEnvString("PORT", "5100"),
		env:  env.GetEnvString("ENV", "development"),
		dbconfig: dbconfig{
			db_url:           env.GetEnvString("DB_ADDR", "postgres://postgres@localhost:5432/pcourierds?sslmode=disable"),
			maxIdleConns:     env.GetEnvInt("SET_MAX_IDLE_CONNS", 10),
			maxOpenConns:     env.GetEnvInt("SET_MAX_OPEN_CONNS", 100),
			connsMaxIdleTime: env.GetEnvTDuration("SET_CONN_MAX_IDLE_TIME", 25*time.Minute),
		},
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	db, err := db.NewPostgresDB(cfg.dbconfig.db_url, cfg.dbconfig.maxIdleConns, cfg.dbconfig.maxOpenConns, cfg.dbconfig.connsMaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Infow("Connected to database successfully")

	app := &application{
		config: cfg,
		logger: logger,
		store:  store.NewStorage(db),
	}

	mux := app.routes()
	logger.Fatal(app.server(mux))
}
