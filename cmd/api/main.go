package main

import (
	"github.com/puremike/pcourierds/internal/env"
	"go.uber.org/zap"
)

var envGet *env.EnvConfig

func main() {

	envGet = env.GetEnv()

	cfg := &config{
		port: envGet.PORT,
		env:  envGet.ENV,
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	app := &app{
		config: cfg,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.start(mux))
}
